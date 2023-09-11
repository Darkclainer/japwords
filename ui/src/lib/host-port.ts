import { err, ok, Result } from 'true-myth/result';

export type HostPort = { host: string; port: string };

// splitHostPort is rewritten version of Go net.SplitHostPort function
export function splitHostPort(hostport: string): Result<HostPort, string> {
  const missingPort = 'missing port in address';
  const tooManyColons = 'too many colons in address';
  const addrErr = (reason: string): Result<HostPort, string> => {
    return err(reason);
  };

  // The port starts after the last colon.
  const i = hostport.lastIndexOf(':');
  if (i < 0) {
    return addrErr(missingPort);
  }
  let j = 0;
  let k = 0;
  let host = '';
  if (hostport[0] == '[') {
    // Expect the first ']' just before the last ':'.
    const end = hostport.indexOf(']');
    if (end < 0) {
      return addrErr("missing ']' in address");
    }
    switch (end + 1) {
      case hostport.length:
        // There can't be a ':' behind the ']' now.
        return addrErr(missingPort);
      case i:
        // The expected result.
        break;
      default:
        // Either ']' isn't followed by a colon, or it is
        // followed by a colon that is not the last one.
        if (hostport[end + 1] == ':') {
          return addrErr(tooManyColons);
        }
        return addrErr(missingPort);
    }
    host = hostport.slice(1, end);
    j = 1;
    k = end + 1;
  } else {
    host = hostport.slice(0, i);
    if (host.indexOf(':') >= 0) {
      return addrErr(tooManyColons);
    }
  }

  if (hostport.indexOf('[', j) >= 0) {
    return addrErr("unexpected '[' in address");
  }
  if (hostport.indexOf(']', k) >= 0) {
    return addrErr("unexpected ']' in address");
  }

  const port = hostport.slice(i + 1);
  return ok({
    host: host,
    port: port,
  });
}
