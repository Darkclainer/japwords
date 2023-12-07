import { describe, expect, test } from '@jest/globals';

import { Lemma, LemmaNoteInfo, Word } from '../api/__generated__/graphql';
import { groupLemmaNotes, isLemmaFitBag, LemmaBag, LemmaBagHeader } from './lemma-bag';

function newWord(init?: (word: Word) => void): Word {
  const word: Word = {
    word: '',
    hiragana: '',
    furigana: [],
    pitchShapes: [],
  };
  if (init) {
    init(word);
  }
  return word;
}

function newBag(init?: (bag: LemmaBag) => void): LemmaBag {
  const bag: LemmaBag = {
    slug: newWord(),
    lemmas: [],
    tags: [],
    forms: [],
    audio: [],
  };
  if (init) {
    init(bag);
  }
  return bag;
}

function newLemma(init?: (lemma: Lemma) => void): Lemma {
  const lemma: Lemma = {
    slug: newWord(),
    definitions: [],
    senseTags: [],
    partsOfSpeech: [],
    tags: [],
    forms: [],
    audio: [],
  };
  if (init) {
    init(lemma);
  }
  return lemma;
}

describe('isLemmaFitBag', () => {
  const testCases: [string, boolean, LemmaBagHeader, Lemma][] = [
    // two empty equal to each other
    ['empty', true, newBag(), newLemma()],

    // this properties doesn't matter
    [
      'definitions',
      true,
      newBag(),
      newLemma((lemma) => {
        lemma.definitions = ['hello'];
      }),
    ],
    [
      'senseTags',
      true,
      newBag(),
      newLemma((lemma) => {
        lemma.senseTags = ['hello'];
      }),
    ],
    [
      'partsOfSpeech',
      true,
      newBag(),
      newLemma((lemma) => {
        lemma.partsOfSpeech = ['hello'];
      }),
    ],

    // this fields matter
    [
      'slug',
      false,
      newBag(),
      newLemma((lemma) => {
        lemma.slug.word = 'hello';
      }),
    ],
    [
      'tags',
      false,
      newBag(),
      newLemma((lemma) => {
        lemma.tags = ['hello'];
      }),
    ],
    [
      'forms',
      false,
      newBag(),
      newLemma((lemma) => {
        lemma.forms = [newWord()];
      }),
    ],
    [
      'audio',
      false,
      newBag(),
      newLemma((lemma) => {
        lemma.audio = [
          {
            source: 'hello',
            type: 'world',
          },
        ];
      }),
    ],

    // check recursion
    [
      'complex',
      true,
      newBag((bag) => {
        bag.slug.hiragana = 'hiragana';
        bag.slug.furigana = [
          {
            kanji: 'h',
            hiragana: 'hi',
          },
          {
            kanji: 'r',
            hiragana: 'ra',
          },
        ];
      }),
      newLemma((lemma) => {
        lemma.slug.hiragana = 'hiragana';
        lemma.slug.furigana = [
          {
            kanji: 'h',
            hiragana: 'hi',
          },
          {
            kanji: 'r',
            hiragana: 'ra',
          },
        ];
      }),
    ],
    [
      'complex',
      false,
      newBag((bag) => {
        bag.slug.hiragana = 'hiragana';
        bag.slug.furigana = [
          {
            kanji: 'h',
            hiragana: 'hi',
          },
          {
            kanji: 'r',
            hiragana: 'NO',
          },
        ];
      }),
      newLemma((lemma) => {
        lemma.slug.hiragana = 'hiragana';
        lemma.slug.furigana = [
          {
            kanji: 'h',
            hiragana: 'hi',
          },
          {
            kanji: 'r',
            hiragana: 'ra',
          },
        ];
      }),
    ],
  ];
  test.each(testCases)('%s: %s', (_, expected, bag, lemma) => {
    const actual = isLemmaFitBag(bag, lemma);
    expect(actual).toEqual(expected);
  });
});

function newLemmaNoteInfo(init?: (lemma: Lemma) => void): LemmaNoteInfo {
  const lemma = newLemma(init);
  return {
    noteID: '',
    lemma: lemma,
  };
}

describe('groupLemmaNotes', () => {
  type BagConstructor = (lemmas: LemmaNoteInfo[]) => LemmaBag;
  const testCases: [string, LemmaNoteInfo[], BagConstructor[]][] = [
    ['empty', [], []],
    [
      'one lemma one bag',
      [newLemmaNoteInfo((l) => (l.slug.hiragana = 'hi'))],
      [
        (lemmas) =>
          newBag((b) => {
            b.slug.hiragana = 'hi';
            b.lemmas = [lemmas[0]];
          }),
      ],
    ],
    [
      'two lemmas one bag',
      [
        newLemmaNoteInfo((l) => {
          l.slug.hiragana = 'hi';
          l.definitions = ['a'];
        }),
        newLemmaNoteInfo((l) => {
          l.slug.hiragana = 'hi';
          l.definitions = ['b'];
        }),
      ],
      [
        (lemmas) =>
          newBag((b) => {
            b.slug.hiragana = 'hi';
            b.lemmas = [lemmas[0], lemmas[1]];
          }),
      ],
    ],
    [
      'two lemmas two bag',
      [
        newLemmaNoteInfo((l) => {
          l.slug.hiragana = 'hi';
          l.definitions = ['a'];
        }),
        newLemmaNoteInfo((l) => {
          l.slug.hiragana = 'ho';
          l.definitions = ['b'];
        }),
      ],
      [
        (lemmas) =>
          newBag((b) => {
            b.slug.hiragana = 'hi';
            b.lemmas = [lemmas[0]];
          }),
        (lemmas) =>
          newBag((b) => {
            b.slug.hiragana = 'ho';
            b.lemmas = [lemmas[1]];
          }),
      ],
    ],
    [
      'alternating lemmas',
      [
        newLemmaNoteInfo((l) => {
          l.slug.hiragana = 'hi';
          l.definitions = ['a'];
        }),
        newLemmaNoteInfo((l) => {
          l.slug.hiragana = 'ho';
          l.definitions = ['b'];
        }),
        newLemmaNoteInfo((l) => {
          l.slug.hiragana = 'hi';
          l.definitions = ['c'];
        }),
      ],
      [
        (lemmas) =>
          newBag((b) => {
            b.slug.hiragana = 'hi';
            b.lemmas = [lemmas[0]];
          }),
        (lemmas) =>
          newBag((b) => {
            b.slug.hiragana = 'ho';
            b.lemmas = [lemmas[1]];
          }),
        (lemmas) =>
          newBag((b) => {
            b.slug.hiragana = 'hi';
            b.lemmas = [lemmas[2]];
          }),
      ],
    ],
  ];
  test.each(testCases)('%s', (_, lemmas, bagConstructors) => {
    const expected = bagConstructors.map((constructor) => constructor(lemmas));
    const actual = groupLemmaNotes(lemmas);
    expect(actual).toEqual(expected);
  });
});
