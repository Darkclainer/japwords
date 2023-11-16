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
        <DeckSelect />
        <SuspenseLoading>
          <NoteMappingSettings />
        </SuspenseLoading>
      </div>
    </HealthStatusPlaceholder>
  );
}

// TODO: probably need suspense too?
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
