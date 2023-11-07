import { useContext } from 'react';
import { Link } from 'react-router-dom';

import StatusIcon, { StatusIconKind } from '../../../components/StatusIcon';
import { HealthStatusContext } from '../../../contexts/health-status';
import { throwErrorHealthStatus } from '../../../model/health-status';
import { DeckSelect } from './deck';
import { NoteSelect } from './note';

export default function AnkiUserSettings() {
  return (
    <HealthStatusPlaceholder>
      <div className="flex flex-col gap-8">
        <DeckSelect />
        <NoteSelect />
      </div>
    </HealthStatusPlaceholder>
  );
}

function HealthStatusPlaceholder({ children }: { children: React.ReactNode }) {
  const props = useHealthStatusProps();
  if (!props) {
    return children;
  }
  return (
    <div className="flex flex-col items-center gap-5 text-xl">
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
