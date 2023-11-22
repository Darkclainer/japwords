import { useMutation, useSuspenseQuery } from '@apollo/client';
import { Label } from '@radix-ui/react-label';
import { useId, useMemo } from 'react';
import { err, ok } from 'true-myth/result';

import { gql } from '../../../api/__generated__';
import { GET_HEALTH_STATUS } from '../../../api/health-status';
import SelectCreate from '../../../components/SelectCreate';
import SuspenseLoading from '../../../components/SuspenseLoading';
import { useToastify } from '../../../hooks/toastify';
import { validateNoteType } from '../../../lib/validate';
import { GET_CURRENT_NOTE } from './api';

export function NoteSelect({ currentNote }: { currentNote: string }) {
  const noteTriggerId = useId();
  return (
    <div className="flex flex-col gap-2.5">
      <Label className="text-2xl" htmlFor={noteTriggerId}>
        Choose a note type:
      </Label>
      <SuspenseLoading>
        <NoteSelectBody triggerId={noteTriggerId} currentNote={currentNote} />
      </SuspenseLoading>
    </div>
  );
}

const GET_ANKI_NOTES = gql(`
  query GetAnkiNotes {
    Anki {
      notes {
        notes
        error {
          __typename
        }
      }
    }
  }
`);

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

function NoteSelectBody({ triggerId, currentNote }: { triggerId: string; currentNote: string }) {
  const [setCurrentNote] = useMutation(SET_CURRENT_NOTE, {
    refetchQueries: [GET_CURRENT_NOTE, GET_HEALTH_STATUS],
    awaitRefetchQueries: true,
  });
  const [createNote] = useMutation(CREATE_DEFAULT_NOTE, {
    refetchQueries: [GET_ANKI_NOTES],
    awaitRefetchQueries: true,
  });
  const { data: notesResp } = useSuspenseQuery(GET_ANKI_NOTES, {
    fetchPolicy: 'network-only',
  });
  const notes = useMemo(() => {
    if (!notesResp.Anki.notes.notes) {
      return null;
    }
    const notes = [...notesResp.Anki.notes.notes];
    return notes.sort().map((item) => {
      return {
        value: item,
      };
    });
  }, [notesResp]);
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
        <p className="text-error-red text-lg">Selected note type does not exists</p>
      )}
    </>
  );
}
