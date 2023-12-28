import { useMutation } from '@apollo/client';
import { Label } from '@radix-ui/react-label';
import { useId, useMemo } from 'react';
import { err, ok } from 'true-myth/result';

import { gql } from '../../../api/__generated__';
import { GET_HEALTH_STATUS } from '../../../api/health-status';
import SelectCreate from '../../../components/SelectCreate';
import SuspenseLoading from '../../../components/SuspenseLoading';
import { useToastify } from '../../../hooks/toastify';
import { validateNoteType } from '../../../lib/validate';
import { GET_ANKI_CONFIG, GET_ANKI_STATE } from './api';

type NoteSelectProps = {
  currentNote: string;
  ankiNotes: string[] | null;
};

export function NoteSelect(props: NoteSelectProps) {
  const noteTriggerId = useId();
  return (
    <div className="flex flex-col gap-2.5">
      <Label className="text-2xl" htmlFor={noteTriggerId}>
        Choose a note type:
      </Label>
      <SuspenseLoading>
        <NoteSelectBody triggerId={noteTriggerId} {...props} />
      </SuspenseLoading>
    </div>
  );
}

const SET_CURRENT_NOTE = gql(`
  mutation SetAnkiConfigCurrentNote($name: String!) {
    setAnkiConfigNote(input: { name: $name }) {
      error {
          message
      }
    }
  }
`);

const CREATE_DEFAULT_NOTE = gql(`
  mutation CreateDefaultAnkiNote($name: String!) {
    createDefaultAnkiNote(input: { name: $name }) {
      ankiError {
        ... on Error {
          message
        }
      }
      error {
        ... on CreateDefaultAnkiNoteAlreadyExists {
          message
        }
        ... on ValidationError {
          message
        }
        ... on Error {
          message
        }
      
      }
    }
  }
`);

function NoteSelectBody({
  triggerId,
  currentNote,
  ankiNotes,
}: { triggerId: string } & NoteSelectProps) {
  const [setCurrentNote] = useMutation(SET_CURRENT_NOTE, {
    refetchQueries: [GET_HEALTH_STATUS, GET_ANKI_CONFIG, GET_ANKI_STATE],
    awaitRefetchQueries: true,
  });
  const [createNote] = useMutation(CREATE_DEFAULT_NOTE, {
    refetchQueries: [GET_HEALTH_STATUS, GET_ANKI_STATE],
    awaitRefetchQueries: true,
  });
  const notes = useMemo(() => {
    if (!ankiNotes) {
      return null;
    }
    const notes = [...ankiNotes];
    return notes.sort().map((item) => {
      return {
        value: item,
      };
    });
  }, [ankiNotes]);
  const toast = useToastify({
    type: 'success',
  });
  if (!notes) {
    // this error handled in parent components
    return null;
  }
  const currentNoteExists = !notes.find((e) => e.value == currentNote);
  return (
    <>
      <SelectCreate
        id={triggerId}
        triggerClassName="max-w-md shrink"
        hasError={currentNoteExists}
        items={notes}
        selectedValue={currentNote}
        onValueChange={async (value: string) => {
          const resp = await setCurrentNote({
            variables: {
              name: value,
            },
          });
          if (!resp.data || resp.data.setAnkiConfigNote.error) {
            toast('Note change failed!', { type: 'error' });
          } else {
            toast('Note successfully changed.');
          }
        }}
        handleCreate={async (value: string) => {
          const resp = await createNote({
            variables: {
              name: value,
            },
          });
          if (!resp.data) {
            toast('Note creation failed!', { type: 'error' });
            return err('request failed');
          }
          if (resp.data.createDefaultAnkiNote.ankiError) {
            toast('Note creation failed! No anki connection', { type: 'error' });
            return err('request failed');
          }
          if (resp.data.createDefaultAnkiNote.error) {
            const error = resp.data.createDefaultAnkiNote.error;
            switch (error.__typename) {
              case 'ValidationError':
                return err(error.message);
              case 'CreateDefaultAnkiNoteAlreadyExists':
                return err('note with specified name already exists');
              default:
                return err('uknown error');
            }
          }
          return ok(value);
        }}
        validateValue={validateNoteType}
        placeholderLabel="Select or create..."
        createLabel="Create new note type"
        createDefaultValue="JapwordsDefaultNoteType"
        dialogTitle="Create note type"
        dialogInputLabel="Input new note type name"
      />
      {currentNoteExists && (
        <p className="text-lg text-error-red">Selected note type does not exists</p>
      )}
    </>
  );
}
