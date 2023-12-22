import { useMutation } from '@apollo/client';
import * as Dialog from '@radix-ui/react-dialog';
import * as Label from '@radix-ui/react-label';
import * as Select from '@radix-ui/react-select';
import { clsx } from 'clsx';
import { FieldArray, Form, Formik, FormikHelpers } from 'formik';
import { Dispatch, SetStateAction, useCallback, useId, useMemo } from 'react';

import { gql } from '../../api/__generated__';
import {
  AddNoteAudioAsset,
  AddNoteAudioAssetInput,
  AddNoteRequest,
  LemmaNoteInfo,
} from '../../api/__generated__/graphql';
import Button, { ButtonVariant } from '../../components/Button';
import { DialogModal, DialogWidth } from '../../components/DialogModal';
import { LoadingIcon } from '../../components/Icons/StatusIcon';
import { SelectItem } from '../../components/Select';
import SelectField from '../../components/SelectField';
import TextField, { TextFieldProps } from '../../components/TextField';
import { apolloErrorToast, ToastFunction } from '../../lib/styled-toast';
import { AddNoteFailedActionTitle } from './model';

export type AddLemmaRequest = {
  note: AddNoteRequest;
  lemma: LemmaNoteInfo;
};

type AddNoteDialogProps = {
  open: boolean;
  setOpen: Dispatch<boolean>;
  addLemmaRequest?: AddLemmaRequest;
  toast: ToastFunction;
  setLemmaNotes: Dispatch<SetStateAction<LemmaNoteInfo[]>>;
};

export default function AddNoteDialog({
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

type AudioPerField = Record<string, AudioVariant[] | undefined>;
type AudioChoices = Record<string, string | undefined>;

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
          <AudioSelectPerField audioPerField={audioPerField} />
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

function AudioSelectPerField({ audioPerField }: { audioPerField: AudioPerField }) {
  return (
    <>
      {Object.entries(audioPerField).map(
        ([field, audios]) => audios && <AudioSelect key={field} field={field} audios={audios} />,
      )}
    </>
  );
}

function AudioSelect({ field, audios }: { field: string; audios: AudioVariant[] }) {
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
