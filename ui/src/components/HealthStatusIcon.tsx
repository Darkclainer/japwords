import { useContext } from 'react';
import { Link } from 'react-router-dom';

import { HealthStatusContext } from '../contexts/health-status';
import { HealthStatus } from '../model/health-status';
import StatusIcon, { StatusIconKind } from './StatusIcon';

function GetIconKind(status: HealthStatus): StatusIconKind {
  switch (status.kind) {
    case 'Ok': {
      switch (status.anki.kind) {
        case 'Ok':
          return 'OK';
        case 'UserError':
        case 'Unauthorized':
        case 'ConnectionError':
        case 'ForbiddenOrigin':
          return 'Warning';
      }
      // linter false positive, but safer anyway
      break;
    }
    case 'Error':
    case 'Disconnected':
      return 'Error';
    case 'Loading':
      return 'Loading';
    default:
      throw 'unreachable';
  }
}

export default function HealthStatusIcon({ size }: { size?: number | string }) {
  const healthStatus = useContext(HealthStatusContext);
  const iconKind = GetIconKind(healthStatus);
  return (
    <Link to="/health-dashboard">
      <StatusIcon kind={iconKind} size={size} />
    </Link>
  );
}
