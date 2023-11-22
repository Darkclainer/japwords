import { useSuspenseQuery } from '@apollo/client';
import { clsx } from 'clsx';
import { useContext } from 'react';
import { Link } from 'react-router-dom';

import StatusIcon, { StatusIconKind } from '../../../components/StatusIcon';
import SuspenseLoading from '../../../components/SuspenseLoading';
import { HealthStatusContext } from '../../../contexts/health-status';
import { AnkiStateOk, throwErrorHealthStatus } from '../../../model/health-status';
import { GET_CURRENT_NOTE } from './api';
import { DeckSelect } from './deck';
import { MappingEdit } from './mapping';
import { NoteSelect } from './note';

export default function AnkiUserSettings() {
  const [ankiState, errorProps] = useAnkiStateOrError();
  if (ankiState) {
    return (
      <div className="flex flex-col gap-8">
        <Notice />
        <DeckSelect />
        <SuspenseLoading>
          <NoteMappingSettings />
        </SuspenseLoading>
        <StatusBox ankiState={ankiState} />
      </div>
    );
  }
  return (
    <div className="flex flex-col items-center gap-5 text-xl my-8">
      <div className="flex flex-row gap-3">
        <StatusIcon kind={errorProps.iconKind} />
        <h1 className="text-3xl">{errorProps.head}</h1>
      </div>
      <div>{errorProps.body}</div>
    </div>
  );
}

function Notice() {
  return (
    <div>
      <p className="text-blue text-lg">
        <span className="font-bold">Note:</span>
        <br />
        After you{' '}
        <Link className="text-dark-blue" to="../connection-settings">
          connected
        </Link>{' '}
        to Anki you should choose deck, note and configure mapping.
        <br />
        Deck is the name of deck where new notes will be added.
        <br />
        Note is the name Anki model that will be used to create new notes. You can ease
        configuration by creating new default note: it is configured to work with default mapping.
        <br />
        Mapping configuration determine how note fields will be filled with information from
        dictionary. For this{' '}
        <a
          className="text-dark-blue"
          href="https://pkg.go.dev/text/template"
          target="_blank"
          rel="noreferrer"
        >
          Go Templates
        </a>{' '}
        were used.
      </p>
    </div>
  );
}

function NoteMappingSettings() {
  const { data: currentNoteResp } = useSuspenseQuery(GET_CURRENT_NOTE, {
    fetchPolicy: 'no-cache',
  });
  const currentNote = currentNoteResp.AnkiConfig.noteType;
  return (
    <>
      <NoteSelect currentNote={currentNote} />
      <MappingEdit currentNote={currentNote} />
    </>
  );
}

type ErrorProps = {
  iconKind: StatusIconKind;
  head: React.ReactNode;
  body?: React.ReactNode;
};

function useAnkiStateOrError(): [null, ErrorProps] | [AnkiStateOk, null] {
  const healthStatus = useContext(HealthStatusContext);
  throwErrorHealthStatus(healthStatus);
  switch (healthStatus.kind) {
    case 'Loading':
      return [
        null,
        {
          iconKind: 'Loading',
          head: 'Loading',
          body: 'Wait a moment',
        },
      ];
    case 'Ok': {
      const ankiState = healthStatus.anki;
      switch (ankiState.kind) {
        case 'Ok':
        case 'UserError':
          return [ankiState, null];
        default:
          return [
            null,
            {
              iconKind: 'Warning',
              head: 'Can not connect to Anki',
              body: (
                <>
                  Configure connection to Anki on{' '}
                  <Link to="../connection-settings" className="text-blue underline">
                    Anki Connect
                  </Link>{' '}
                  page
                </>
              ),
            },
          ];
      }
    }
    default: {
      const _exhaustiveCheck: never = healthStatus;
      return _exhaustiveCheck;
    }
  }
}

function StatusBox({ ankiState }: { ankiState: AnkiStateOk }) {
  type error = {
    key: string;
    msg: string;
  };
  const errors: error[] = [];
  if (!ankiState.deckExists) {
    errors.push({
      key: 'nodeck',
      msg: "Selected deck doesn't exists.",
    });
  }
  if (!ankiState.noteTypeExists) {
    errors.push({
      key: 'nonote',
      msg: "Selected note type doesn't exists.",
    });
  }
  if (!ankiState.noteHasAllFields) {
    errors.push({
      key: 'invalidfields',
      msg: 'Mapping has fields that note type has not.',
    });
  }
  return (
    <div className="flex flex-row justify-start items-start gap-2 basis-16">
      <StatusIcon size="2.5rem" kind={errors.length == 0 ? 'OK' : 'Error'} />
      <div className="text-2xl">
        <h1
          className={clsx(
            'text-bold leading-10',
            errors.length == 0 ? 'text-green' : 'text-error-red',
          )}
        >
          {errors.length == 0 ? 'OK' : 'Error'}
        </h1>
        {errors.length == 0 && 'All is configured, you can add words!'}
        {errors.map((err) => (
          <p key={err.key}>{err.msg}</p>
        ))}
      </div>
    </div>
  );
}
