import { describe, expect, test } from '@jest/globals';

import { splitHostPort } from './host-port';

describe('no error', () => {
  const testCases = [
    // Host name
    ['localhost:http', 'localhost', 'http'],
    ['localhost:80', 'localhost', '80'],

    // Go-specific host name with zone identifier
    ['localhost%lo0:http', 'localhost%lo0', 'http'],
    ['localhost%lo0:80', 'localhost%lo0', '80'],
    ['[localhost%lo0]:http', 'localhost%lo0', 'http'], // Go 1 behavior
    ['[localhost%lo0]:80', 'localhost%lo0', '80'], // Go 1 behavior

    // IP literal
    ['127.0.0.1:http', '127.0.0.1', 'http'],
    ['127.0.0.1:80', '127.0.0.1', '80'],
    ['[::1]:http', '::1', 'http'],
    ['[::1]:80', '::1', '80'],

    // IP literal with zone identifier
    ['[::1%lo0]:http', '::1%lo0', 'http'],
    ['[::1%lo0]:80', '::1%lo0', '80'],

    // Go-specific wildcard for host name
    [':http', '', 'http'], // Go 1 behavior
    [':80', '', '80'], // Go 1 behavior

    // Go-specific wildcard for service name or transport port number
    ['golang.org:', 'golang.org', ''], // Go 1 behavior
    ['127.0.0.1:', '127.0.0.1', ''], // Go 1 behavior
    ['[::1]:', '::1', ''], // Go 1 behavior

    // Opaque service name
    ['golang.org:https%foo', 'golang.org', 'https%foo'], // Go 1 behavior
  ];
  test.each(testCases)('hostport: %s', (src, host, port) => {
    const actual = splitHostPort(src);
    expect(actual.isOk).toBeTruthy();
    const actualOk = actual.unwrapOr({ host: '', port: '' });
    expect(actualOk.host).toEqual(host);
    expect(actualOk.port).toEqual(port);
  });
});

describe('error', () => {
  const testCases = [
    ['golang.org', 'missing port in address'],
    ['127.0.0.1', 'missing port in address'],
    ['[::1]', 'missing port in address'],
    ['[fe80::1%lo0]', 'missing port in address'],
    ['[localhost%lo0]', 'missing port in address'],
    ['localhost%lo0', 'missing port in address'],

    ['::1', 'too many colons in address'],
    ['fe80::1%lo0', 'too many colons in address'],
    ['fe80::1%lo0:80', 'too many colons in address'],

    // Test cases that didn't fail in Go 1

    ['[foo:bar]', 'missing port in address'],
    ['[foo:bar]baz', 'missing port in address'],
    ['[foo]bar:baz', 'missing port in address'],

    ['[foo]:[bar]:baz', 'too many colons in address'],

    ['[foo]:[bar]baz', "unexpected '[' in address"],
    ['foo[bar]:baz', "unexpected '[' in address"],

    ['foo]bar:baz', "unexpected ']' in address"],
  ];
  test.each(testCases)('hostport: %s', (src, error) => {
    const actual = splitHostPort(src);
    if (actual.variant == 'Ok') {
      throw 'expected error';
    }
    expect(actual.error).toEqual(error);
  });
});
