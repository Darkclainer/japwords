import { useLazyQuery, useMutation } from '@apollo/client';
import * as Dialog from '@radix-ui/react-dialog';
import * as Label from '@radix-ui/react-label';
import { clsx } from 'clsx';
import { FieldArray, Form, Formik, FormikHelpers } from 'formik';
import { Dispatch, useCallback, useId, useMemo, useRef, useState } from 'react';

import { gql } from '../api/__generated__';
import { AddNoteRequest, Lemma } from '../api/__generated__/graphql';
import { useToastify } from '../hooks/toastify';
import { apolloErrorToast, ToastFunction } from '../lib/styled-toast';
import Button, { ButtonVariant } from './Button';
import { DialogModal, DialogWidth } from './DialogModal';
import { LoadingIcon } from './Icons/StatusIcon';
import LemmaCard from './LemmaCard';
import TextField, { TextFieldProps } from './TextField';

const AddNoteFailedActionTitle = 'Add note failed';

const PREPARE_PROJECTED_LEMMA = gql(`
query PrepareProjectedLemma($lemma: ProjectedLemmaInput) {
  PrepareProjectedLemma(lemma: $lemma) {
    request {
      fields {
        name
        value
      }
      tags
      audioURL
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

export default function LemmaList({ lemmas }: { lemmas: Array<Lemma> }) {
  const [prepareLemma] = useLazyQuery(PREPARE_PROJECTED_LEMMA, {
    fetchPolicy: 'network-only',
  });
  const [openAddNote, setOpenAddNote] = useState(false);
  const [noteRequest, setNoteRequest] = useState<AddNoteRequest>();
  const abortReset = useAbortRef();
  const toast = useToastify({
    autoClose: 2000,
    type: 'info',
  });
  const previewLemma = useCallback(
    async (lemma: Lemma, senseIndex: number) => {
      const currentAbortController = abortReset();
      setNoteRequest(undefined);
      setOpenAddNote(true);
      const definition = lemma.senses[senseIndex];
      const now = new Date().valueOf();
      const { data, error } = await prepareLemma({
        variables: {
          lemma: {
            slug: lemma.slug,
            forms: lemma.forms,
            tags: lemma.tags,
            definitions: definition.definition,
            senseTags: definition.tags,
            partsOfSpeech: definition.partOfSpeech,
            audio: lemma.audio,
          },
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
      const ankiError = data?.PrepareProjectedLemma.ankiError;
      if (ankiError) {
        setOpenAddNote(false);
        toast(`${AddNoteFailedActionTitle}: problems with anki-connect`, {
          type: 'error',
        });
        return;
      }
      const userError = data?.PrepareProjectedLemma.error;
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
        setNoteRequest(data?.PrepareProjectedLemma.request ?? undefined);
      }
    },
    [setOpenAddNote, toast, prepareLemma],
  );

  return (
    <>
      {lemmas.map((lemma, index) => (
        <LemmaCard
          key={lemma.slug.word + '-' + index}
          lemma={lemma}
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
        noteRequest={noteRequest}
        toast={toast}
      />
    </>
  );
}

type AddNoteDialogProps = {
  open: boolean;
  setOpen: Dispatch<boolean>;
  noteRequest?: AddNoteRequest;
  toast: ToastFunction;
};

function AddNoteDialog({ open, setOpen, noteRequest, toast }: AddNoteDialogProps) {
  return (
    <DialogModal open={open} onOpenChange={setOpen} widthVariant={DialogWidth.Large}>
      <Dialog.Title className="mb-2.5 text-2xl font-bold text-blue">Add note to Anki</Dialog.Title>

      {noteRequest === undefined ? (
        <div className="flex p-8 text-2xl gap-5 justify-center place-items-center">
          <div>
            <LoadingIcon className="w-8 h-8" />
          </div>
          <div>Loading...</div>
        </div>
      ) : (
        <AddNoteForm noteRequest={noteRequest} setOpen={setOpen} toast={toast} />
      )}
    </DialogModal>
  );
}
const ADD_ANKI_NOTE = gql(`
mutation AddAnkiNote($note: AddNoteRequestInput!) {
  addAnkiNote(request: $note) {
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
  noteRequest: AddNoteRequest;
  toast: ToastFunction;
  setOpen: Dispatch<boolean>;
};

type Note = {
  fields: NoteField[];
  audioURL: string;
  tags: string[];
};

type NoteField = {
  name: string;
  value: string;
};

function AddNoteForm({ noteRequest, setOpen, toast }: AddNoteFormProps) {
  const note: Note = useMemo(() => {
    return {
      fields: noteRequest.fields.map((field) => ({
        name: field.name,
        value: field.value,
      })),
      audioURL: '',
      tags: [],
    };
  }, [noteRequest]);
  const [addNote] = useMutation(ADD_ANKI_NOTE, {
    // we passe empty onError because in this case apollo will not throw in case if network error for some reason
    onError: () => {},
  });
  const handleSubmit = useCallback(async (values: Note, helpers: FormikHelpers<Note>) => {
    console.log(values);
    const { data, errors } = await addNote({
      variables: {
        note: note,
      },
    });
    if (errors) {
      apolloErrorToast(errors, `${AddNoteFailedActionTitle}: `, { toast: toast });
      return;
    }
    if (data?.addAnkiNote.ankiError) {
      toast(`${AddNoteFailedActionTitle}: problems with anki-connect`, {
        type: 'error',
      });
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
    toast('Note added to Anki');
    setOpen(false);
  }, []);
  return (
    <Formik<Note> initialValues={note} onSubmit={handleSubmit}>
      {(props) => (
        <Form className="flex flex-col gap-2.5">
          <FieldArray
            name="fields"
            render={() =>
              noteRequest.fields.map((field, index) => (
                <FormTextField key={field.name} name={`fields.${index}.value`} label={field.name} />
              ))
            }
          />
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
