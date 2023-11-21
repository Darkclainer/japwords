import { useContext } from 'react';
import { Link } from 'react-router-dom';

import StatusIcon, { StatusIconKind } from '../../../components/StatusIcon';
import { HealthStatusContext } from '../../../contexts/health-status';
import { throwErrorHealthStatus } from '../../../model/health-status';
import { DeckSelect } from './deck';
import { MappingEdit } from './mapping';
import { NoteSelect } from './note';
import { useSuspenseQuery } from '@apollo/client';
import { GET_CURRENT_NOTE } from './api';
import SuspenseLoading from '../../../components/SuspenseLoading';

export default function AnkiUserSettings() {
  return (
    <HealthStatusPlaceholder>
      <div className="flex flex-col gap-8">
        <Notice />
        <DeckSelect />
        <SuspenseLoading>
          <NoteMappingSettings />
        </SuspenseLoading>
      </div>
    </HealthStatusPlaceholder>
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
    fetchPolicy: 'network-only',
  });
  const currentNote = currentNoteResp.AnkiConfig.noteType;
  return (
    <>
      <NoteSelect currentNote={currentNote} />
      <MappingEdit currentNote={currentNote} />
    </>
  );
}

function HealthStatusPlaceholder({ children }: { children: React.ReactNode }) {
  const props = useHealthStatusProps();
  if (!props) {
    return children;
  }
  return (
    <div className="flex flex-col items-center gap-5 text-xl my-8">
      <div className="flex flex-row gap-3">
        <StatusIcon kind={props.iconKind} />
        <h1 className="text-3xl">{props.head}</h1>
      </div>
      <div>{props.body}</div>
    </div>
  );
}

function useHealthStatusProps(): {
  iconKind: StatusIconKind;
  head: React.ReactNode;
  body?: React.ReactNode;
} | null {
  const healthStatus = useContext(HealthStatusContext);
  throwErrorHealthStatus(healthStatus);
  switch (healthStatus.kind) {
    case 'Loading':
      return {
        iconKind: 'Loading',
        head: 'Loading',
        body: 'Wait a moment',
      };
    case 'Ok': {
      const ankiState = healthStatus.anki;
      switch (ankiState.kind) {
        case 'Ok':
        case 'UserError':
          return null;
        default:
          return {
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
          };
      }
    }
    default: {
      const _exhaustiveCheck: never = healthStatus;
      return _exhaustiveCheck;
    }
  }
}
