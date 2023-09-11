import { splitHostPort } from './host-port';

const hostnameRfc1123 =
  /^([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62}){1}(\.[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})*?$/;

// validateHostPort accepts strings that consists of hostname compaining RFC1123 and port separated by colon.
export function validateHostPort(src: string): string | null {
  const hostPort = splitHostPort(src);
  if (hostPort.isErr) {
    return hostPort.error;
  }
  const port = Number(hostPort.value.port);
  if (isNaN(port) || !Number.isInteger(port)) {
    return 'port is not an integer';
  }
  if (port < 1 || port > 65535) {
    return 'port value should be above 0 and less than 65536';
  }
  if (!hostnameRfc1123.test(hostPort.value.host)) {
    return 'hostname is not valid';
  }
  return null;
}

export function validateDeck(name: string): string | null {
  if (!validateDeckNoteType(name)) {
    return 'deck name is not valid';
  }
  return null;
}

export function validateNoteType(name: string): string | null {
  if (!validateDeckNoteType(name)) {
    return 'note type name is not valid';
  }
  return null;
}

const deckNoteType = /(^[^ \t\n\v"][^"\n]*[^ \t\n\v"]$)|(^[^ \t\n\v"]$)/;

function validateDeckNoteType(name: string): boolean {
  return deckNoteType.test(name);
}
