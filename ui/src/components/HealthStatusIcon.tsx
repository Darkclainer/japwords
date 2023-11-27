import { useContext } from 'react';
import { Link } from 'react-router-dom';

import { HealthStatusContext } from '../contexts/health-status';
import { HealthStatus } from '../model/health-status';
import { IconProps } from './Icons/IconProps';
import StatusIcon, { StatusIconKind } from './Icons/StatusIcon';

function GetIconKind(status: HealthStatus): StatusIconKind {
  switch (status.kind) {
    case 'Ok': {
      switch (status.anki.kind) {
        case 'Ok':
          return 'OK';
        case 'UserError':
        case 'InvalidAPIKey':
        case 'CollectionUnavailable':
        case 'ConnectionError':
        case 'ForbiddenOrigin':
          return 'Warning';
        case 'UnknownError':
          return 'Error';
        default: {
          const _exhaustiveCheck: never = status.anki;
          return _exhaustiveCheck;
        }
      }
    }
    case 'Error':
    case 'Disconnected':
      return 'Error';
    case 'Loading':
      return 'Loading';
    default: {
      const _exhaustiveCheck: never = status;
      return _exhaustiveCheck;
    }
  }
}

export default function HealthStatusIcon(props: IconProps) {
  const healthStatus = useContext(HealthStatusContext);
  const iconKind = GetIconKind(healthStatus);
  return (
    <Link to="/health-dashboard">
      <StatusIcon kind={iconKind} {...props} />
    </Link>
  );
}
