import { describe, expect, test } from '@jest/globals';

import { validateDeck, validateHostPort, validateNoteType } from './validate';

describe('validateHostPort', () => {
  const testCases = [
    // positive
    { src: '127.0.0.1:8765', expected: null },
    { src: '127.0.0.1:1', expected: null },
    { src: '127.0.0.1:65535', expected: null },
    { src: 'example.com:8765', expected: null },
    { src: 'localhost:8765', expected: null },
    { src: 'localhost:000008765', expected: null },
    // negative
    { src: '', expected: 'missing port in address' },
    { src: '127.0.0.1', expected: 'missing port in address' },
    { src: ':8765', expected: 'hostname is not valid' },
    { src: '127.0.0.1:0', expected: 'port value should be above 0 and less than 65536' },
    { src: '127.0.0.1:65536', expected: 'port value should be above 0 and less than 65536' },
    { src: '127.0.0.1:1.5', expected: 'port is not an integer' },
    { src: '127.0.0.1:a', expected: 'port is not an integer' },
    { src: '127.0.0.1:-80', expected: 'port value should be above 0 and less than 65536' },
    { src: '127.0.0.1::80', expected: 'too many colons in address' },
    { src: 'example..com:80', expected: 'hostname is not valid' },
    { src: 'example.com.:80', expected: 'hostname is not valid' },
  ];
  test.each(testCases)('src: $src', ({ src, expected }) => {
    const actual = validateHostPort(src);
    expect(actual).toEqual(expected);
  });
});

describe('validateDeck', () => {
  const errorMessage = 'deck name is not valid';
  const testCases = [
    // positive
    { src: 'a', expected: null },
    { src: 'ab', expected: null },
    { src: 'a  b', expected: null },
    { src: 'a#!b', expected: null },
    { src: 'a{b', expected: null },
    { src: 'a}b', expected: null },
    { src: 'a:b', expected: null },
    // negative
    { src: '', expected: errorMessage },
    { src: ' ', expected: errorMessage },
    { src: '  ', expected: errorMessage },
    { src: ' a', expected: errorMessage },
    { src: 'a ', expected: errorMessage },
    { src: ' a ', expected: errorMessage },
    { src: ' :', expected: errorMessage },
    { src: ': ', expected: errorMessage },
    { src: 'a"b', expected: errorMessage },
    { src: 'a: ', expected: errorMessage },
    { src: ' :a', expected: errorMessage },
    { src: 'a\nb', expected: errorMessage },
  ];
  test.each(testCases)('src: $src', ({ src, expected }) => {
    const actual = validateDeck(src);
    expect(actual).toEqual(expected);
  });
});

describe('validateNoteType', () => {
  const errorMessage = 'note type name is not valid';
  const testCases = [
    // positive
    { src: 'a', expected: null },
    { src: 'ab', expected: null },
    { src: 'a  b', expected: null },
    { src: 'a#!b', expected: null },
    { src: 'a{b', expected: null },
    { src: 'a}b', expected: null },
    { src: 'a:b', expected: null },
    // negative
    { src: '', expected: errorMessage },
    { src: ' ', expected: errorMessage },
    { src: '  ', expected: errorMessage },
    { src: ' a', expected: errorMessage },
    { src: 'a ', expected: errorMessage },
    { src: ' a ', expected: errorMessage },
    { src: ' :', expected: errorMessage },
    { src: ': ', expected: errorMessage },
    { src: 'a"b', expected: errorMessage },
    { src: 'a: ', expected: errorMessage },
    { src: ' :a', expected: errorMessage },
    { src: 'a\nb', expected: errorMessage },
  ];
  test.each(testCases)('src: $src', ({ src, expected }) => {
    const actual = validateNoteType(src);
    expect(actual).toEqual(expected);
  });
});
