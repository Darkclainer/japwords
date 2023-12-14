import { useLazyQuery, useMutation } from '@apollo/client';
import * as Dialog from '@radix-ui/react-dialog';
import * as Label from '@radix-ui/react-label';
import * as Select from '@radix-ui/react-select';
import { clsx } from 'clsx';
import { FieldArray, Form, Formik, FormikHelpers } from 'formik';
import { Dispatch, SetStateAction, useCallback, useId, useMemo, useRef, useState } from 'react';

import { gql } from '../api/__generated__';
import {
  AddNoteAudioAsset,
  AddNoteAudioAssetInput,
  AddNoteRequest,
  LemmaNoteInfo,
} from '../api/__generated__/graphql';
import { useToastify } from '../hooks/toastify';
import { groupLemmaNotes } from '../lib/lemma-bag';
import { apolloErrorToast, ToastFunction } from '../lib/styled-toast';
import Button, { ButtonVariant } from './Button';
import { DialogModal, DialogWidth } from './DialogModal';
import { LoadingIcon } from './Icons/StatusIcon';
import LemmaCard from './LemmaCard';
import { SelectItem } from './Select';
import SelectField from './SelectField';
import TextField, { TextFieldProps } from './TextField';

const AddNoteFailedActionTitle = 'Add note failed';

const PREPARE_LEMMA = gql(`
query PrepareLemma($lemma: LemmaInput) {
  PrepareLemma(lemma: $lemma) {
    request {
      fields {
        name
        value
      }
      tags
      audioAssets {
        field
        filename
        url
        data
      }
    }
    error {
      ... on AnkiIncompleteConfiguration {
        message
      }
      ... on Error {
        message
      }
    }
    ankiError {
      __typename
      ... on Error {
        message
      }
    }
  }
}
`);

type AddLemmaRequest = {
  note: AddNoteRequest;
  lemma: LemmaNoteInfo;
};

export default function LemmaList({
  lemmaNotes: initialLemmaNotes,
}: {
  lemmaNotes: LemmaNoteInfo[];
}) {
  // more propper way will be to modify cache, but because lemmas now returned from react-router-dom it's not possible
  const [lemmaNotes, setLemmaNotes] = useState(initialLemmaNotes);
  const lemmaBags = useMemo(() => groupLemmaNotes(lemmaNotes), [lemmaNotes]);
  const [prepareLemma] = useLazyQuery(PREPARE_LEMMA, {
    fetchPolicy: 'network-only',
  });
  const [openAddNote, setOpenAddNote] = useState(false);
  const [addLemmaRequest, setAddLemmaRequest] = useState<AddLemmaRequest>();
  const abortReset = useAbortRef();
  const toast = useToastify({
    autoClose: 2000,
    type: 'info',
  });
  const previewLemma = useCallback(
    async (lemmaNote: LemmaNoteInfo) => {
      const currentAbortController = abortReset();
      setAddLemmaRequest(undefined);
      setOpenAddNote(true);
      const now = new Date().valueOf();
      const { data, error } = await prepareLemma({
        variables: {
          lemma: lemmaNote.lemma,
        },
        context: {
          fetchOptions: {
            signal: currentAbortController.signal,
          },
        },
      });
      // if request is to fast, display loading for minimum 400ms
      // Better, to show window if request is longer then 200ms for example?
      const waitTime = Math.max(0, -(new Date().valueOf() - now) + 400);
      if (waitTime > 0) {
        await new Promise((resolve) => setTimeout(resolve, waitTime));
      }
      if (error) {
        setOpenAddNote(false);
        apolloErrorToast(error, `${AddNoteFailedActionTitle}: `, {
          toast,
        });
        return;
      }
      const ankiError = data?.PrepareLemma.ankiError;
      if (ankiError) {
        setOpenAddNote(false);
        toast(`${AddNoteFailedActionTitle}: problems with anki-connect`, {
          type: 'error',
        });
        return;
      }
      const userError = data?.PrepareLemma.error;
      if (userError) {
        setOpenAddNote(false);
        let message = `${AddNoteFailedActionTitle}: `;
        if (userError.__typename == 'AnkiIncompleteConfiguration') {
          message += 'anki configuration unfinished';
        } else {
          message += 'uknown reason';
        }
        toast(message, {
          type: 'error',
        });
        return;
      }
      if (!currentAbortController.signal.aborted) {
        if (!data?.PrepareLemma.request) {
          throw 'unreachable';
        }
        setAddLemmaRequest({
          note: data.PrepareLemma.request,
          lemma: lemmaNote,
        });
      }
    },
    [setOpenAddNote, toast, prepareLemma],
  );

  return (
    <>
      {lemmaBags.map((lemmaBag, index) => (
        <LemmaCard
          key={lemmaBag.slug.word + '-' + index}
          lemmaBag={lemmaBag}
          toast={toast}
          previewLemma={previewLemma}
        />
      ))}
      <AddNoteDialog
        open={openAddNote}
        setOpen={(value: boolean) => {
          if (!value) {
            abortReset();
          }
          return setOpenAddNote(value);
        }}
        addLemmaRequest={addLemmaRequest}
        toast={toast}
        setLemmaNotes={setLemmaNotes}
      />
    </>
  );
}

type AddNoteDialogProps = {
  open: boolean;
  setOpen: Dispatch<boolean>;
  addLemmaRequest?: AddLemmaRequest;
  toast: ToastFunction;
  setLemmaNotes: Dispatch<SetStateAction<LemmaNoteInfo[]>>;
};

function AddNoteDialog({
  open,
  setOpen,
  addLemmaRequest,
  setLemmaNotes,
  toast,
}: AddNoteDialogProps) {
  return (
    <DialogModal open={open} onOpenChange={setOpen} widthVariant={DialogWidth.Large}>
      <Dialog.Title className="mb-2.5 text-2xl font-bold text-blue">Add note to Anki</Dialog.Title>

      {addLemmaRequest === undefined ? (
        <div className="flex p-8 text-2xl gap-5 justify-center place-items-center">
          <div>
            <LoadingIcon className="w-8 h-8" />
          </div>
          <div>Loading...</div>
        </div>
      ) : (
        <AddNoteForm
          addLemmaRequest={addLemmaRequest}
          setLemmaNotes={setLemmaNotes}
          setOpen={setOpen}
          toast={toast}
        />
      )}
    </DialogModal>
  );
}
const ADD_ANKI_NOTE = gql(`
mutation AddAnkiNote($note: AddNoteRequestInput!) {
  addAnkiNote(request: $note) {
    noteID
    error {
      ... on AnkiIncompleteConfiguration {
        message
      }
      ... on AnkiAddNoteDuplicateFound {
        message
      }
    }
    ankiError {
      ... on Error {
        message
      }
    }
  }
}
`);

type AddNoteFormProps = {
  addLemmaRequest: AddLemmaRequest;
  toast: ToastFunction;
  setOpen: Dispatch<boolean>;
  setLemmaNotes: Dispatch<SetStateAction<LemmaNoteInfo[]>>;
};

type Note = {
  fields: NoteField[];
  tags: string[];
  audioChoices: AudioChoices;
};

type AudioVariant = {
  filename: string;
  data: string;
  url: string;
};

type NoteField = {
  name: string;
  value: string;
};

function AddNoteForm({ addLemmaRequest, setOpen, setLemmaNotes, toast }: AddNoteFormProps) {
  const [note, audioPerField] = useMemo(() => {
    const [audioPerField, audioChoices] = groupAudioAssetsPerField(
      addLemmaRequest.note.audioAssets,
    );
    return [
      {
        fields: addLemmaRequest.note.fields.map((field) => ({
          name: field.name,
          value: field.value,
        })),
        tags: [],
        audioChoices: audioChoices,
      },
      audioPerField,
    ];
  }, [addLemmaRequest]);
  const [addNote] = useMutation(ADD_ANKI_NOTE, {
    // we passe empty onError because in this case apollo will not throw in case if network error for some reason
    onError: () => {},
  });
  const handleSubmit = useCallback(
    async (values: Note, helpers: FormikHelpers<Note>) => {
      const audioAssets: AddNoteAudioAssetInput[] = Object.entries(values.audioChoices)
        .map(([field, choice]) => {
          if (choice === undefined) {
            return undefined;
          }
          const asset = audioPerField[field]?.[Number(choice)];
          if (asset === undefined) {
            return undefined;
          }
          return {
            field: field,
            url: asset.url,
            data: asset.data,
            filename: asset.filename,
          };
        })
        .filter(<T,>(item: T | undefined): item is T => item !== undefined);
      const { data, errors } = await addNote({
        variables: {
          note: {
            fields: values.fields,
            tags: values.tags,
            audioAssets: audioAssets,
          },
        },
      });
      if (errors) {
        apolloErrorToast(errors, `${AddNoteFailedActionTitle}: `, { toast: toast });
        return;
      }
      if (data && data.addAnkiNote.noteID !== '') {
        toast('Note added to Anki', { type: 'success' });
        setOpen(false);
        setLemmaNotes((lemmas) =>
          lemmas.map((lemma) => {
            if (lemma === addLemmaRequest.lemma) {
              return {
                lemma: lemma.lemma,
                // must be defined if no errors
                noteID: data.addAnkiNote.noteID,
              };
            }
            {
              return lemma;
            }
          }),
        );
      }
      if (data?.addAnkiNote.ankiError) {
        toast(`${AddNoteFailedActionTitle}: problems with anki-connect`, {
          type: 'error',
        });
        return;
      }
      const userError = data?.addAnkiNote.error;
      if (userError) {
        switch (userError.__typename) {
          case 'AnkiAddNoteDuplicateFound':
            helpers.setErrors({
              fields: [{ value: 'duplicated found' }],
            });
            toast(`${AddNoteFailedActionTitle}: similar note already exists in Anki`, {
              type: 'error',
            });
            return;
          case 'AnkiIncompleteConfiguration':
            toast(`${AddNoteFailedActionTitle}: anki configuration is incomplete`, {
              type: 'error',
            });
            return;
          case undefined:
            throw 'unreachable';
          default: {
            const _exhaustiveCheck: never = userError;
            return _exhaustiveCheck;
          }
        }
      }
    },
    [addLemmaRequest, setLemmaNotes, setOpen, toast],
  );
  return (
    <Formik<Note> initialValues={note} onSubmit={handleSubmit}>
      {(props) => (
        <Form className="flex flex-col gap-2.5">
          <FieldArray
            name="fields"
            render={() =>
              addLemmaRequest.note.fields.map((field, index) => (
                <FormTextField key={field.name} name={`fields.${index}.value`} label={field.name} />
              ))
            }
          />
          <AudioChoiceGroups audioPerField={audioPerField} />
          <div className="flex flex-row gap-8 mt-4">
            <Button type="submit" className="basis-52" disabled={props.isSubmitting}>
              Add Note
            </Button>
            <Button
              type="button"
              onClick={props.handleReset}
              className={clsx('basis-52')}
              variant={ButtonVariant.Dangerous}
              disabled={props.isSubmitting}
            >
              Reset
            </Button>
          </div>
        </Form>
      )}
    </Formik>
  );
}

function AudioChoiceGroups({ audioPerField }: { audioPerField: AudioPerField }) {
  return (
    <>
      {Object.entries(audioPerField).map(
        ([field, audios]) => audios && <AudioChoices key={field} field={field} audios={audios} />,
      )}
    </>
  );
}

function AudioChoices({ field, audios }: { field: string; audios: AudioVariant[] }) {
  const groupId = useId();
  const audioResources = useMemo(() => {
    return audios.map((audio) => {
      const audioResource = new Audio(audio.url);
      audioResource.preload = 'none';
      return audioResource;
    });
  }, [audios]);
  const playAudio = (strIndex: string) => {
    const audio = audioResources[Number(strIndex)];
    if (audio !== undefined) {
      audio.play();
    }
  };
  return (
    <div className={clsx('flex flex-col gap-1.5 text-xl')}>
      <Label.Root htmlFor={groupId}>Audio asset for field: {field}</Label.Root>
      <SelectField
        id={groupId}
        name={`audioChoices.${field}`}
        placeholderLabel={'hello'}
        onValueChange={playAudio}
      >
        <SelectItem value="none">
          <Select.ItemText>None</Select.ItemText>
        </SelectItem>
        {audios.map((audio, index) => (
          <SelectItem key={index} value={index.toString()}>
            <Select.ItemText>{audio.filename}</Select.ItemText>
          </SelectItem>
        ))}
      </SelectField>
    </div>
  );
}

type AudioPerField = Record<string, AudioVariant[] | undefined>;
type AudioChoices = Record<string, string | undefined>;

// returns possible audio assets per field and default choices for these assets (first one if exists)
function groupAudioAssetsPerField(audios: AddNoteAudioAsset[]): [AudioPerField, AudioChoices] {
  const audioPerField: AudioPerField = {};
  audios.forEach((audio) => {
    let variants = audioPerField[audio.field];
    if (variants === undefined) {
      variants = [];
      audioPerField[audio.field] = variants;
    }
    variants.push({
      filename: audio.filename,
      url: audio.url,
      data: audio.data,
    });
  });
  const audioChoices: AudioChoices = {};
  Object.entries(audioPerField).forEach(([field, assets]) => {
    audioChoices[field] = assets && assets.length > 0 ? '0' : '';
  });
  return [audioPerField, audioChoices];
}

type FormTextFieldProps = TextFieldProps & { label: string };

function FormTextField({ name, label, ...rest }: FormTextFieldProps) {
  const inputId = useId();
  return (
    <div className={clsx('flex flex-col gap-1.5 text-xl')}>
      {label && <Label.Root htmlFor={inputId}>{label}:</Label.Root>}
      <div className="flex  flex-col  gap-2.5">
        <TextField id={inputId} name={name} {...rest} />
      </div>
    </div>
  );
}

function useAbortRef() {
  const abortRef = useRef<AbortController>();
  const abortReset = useCallback(() => {
    const previousController = abortRef.current;
    if (previousController) {
      previousController.abort();
    }
    const newController = new AbortController();
    abortRef.current = newController;
    return newController;
  }, []);
  return abortReset;
}
