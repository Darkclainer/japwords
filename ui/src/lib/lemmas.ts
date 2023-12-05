import { isEqual } from 'lodash';

import { Audio, Lemma, LemmaNoteInfo, Word } from '../api/__generated__/graphql';

export type LemmaBagHeader = {
  slug: Word;
  tags: string[];
  forms: Word[];
  audio: Audio[];
};

export type LemmaBag = {
  lemmas: LemmaNoteInfo[];
} & LemmaBagHeader;

export function groupLemmaNotes(lemmas: LemmaNoteInfo[]): LemmaBag[] {
  const bags: LemmaBag[] = [];
  let bag: LemmaBag | null = null;
  for (let i = 0; i < lemmas.length; i++) {
    const lemma = lemmas[i];
    if (bag === null) {
      bag = newBagFromLemma(lemma);
      continue;
    }
    if (isLemmaFitBag(bag, lemma.lemma)) {
      bag.lemmas.push(lemma);
    } else {
      bags.push(bag);
      bag = newBagFromLemma(lemma);
    }
  }
  if (bag !== null) {
    bags.push(bag);
  }
  return bags;
}

function newBagFromLemma(lemma: LemmaNoteInfo): LemmaBag {
  return {
    lemmas: [lemma],
    slug: lemma.lemma.slug,
    tags: lemma.lemma.tags,
    forms: lemma.lemma.forms,
    audio: lemma.lemma.audio,
  };
}

function isLemmaFitBag(bag: LemmaBagHeader, lemma: Lemma): boolean {
  const properties: Array<keyof LemmaBagHeader> = ['slug', 'tags', 'forms', 'audio'];
  return properties.every((property) => {
    console.log(
      'bag',
      bag,
      'lemma',
      lemma,
      'property',
      property,
      'result',
      isEqual(bag[property], lemma[property]),
    );
    return isEqual(bag[property], lemma[property]);
  });
}
