import { clsx } from 'clsx';
import { ReactNode, useRef } from 'react';

import {
  AccentDirection,
  Audio,
  Furigana,
  LemmaNoteInfo,
  Word,
} from '../../api/__generated__/graphql';
import PlayIcon from '../../components/Icons/PlayIcon';
import { LemmaBag } from '../../lib/lemma-bag';
import { ToastFunction } from '../../lib/styled-toast';

export type LemmaCardProps = {
  lemmaBag: LemmaBag;
  toast: ToastFunction;
  previewLemma: (lemma: LemmaNoteInfo) => Promise<void>;
};

export default function LemmaCard({ lemmaBag, toast, previewLemma }: LemmaCardProps) {
  return (
    <div className="my-4 rounded-md bg-gray px-10 py-2 shadow-md">
      <div className="flex flex-col divide-y divide-blue">
        <LemmaCardItem render={true}>
          <LemmaTitle word={lemmaBag.slug} toast={toast} />
        </LemmaCardItem>
        <LemmaCardItem render={lemmaBag.slug.hiragana !== '' || lemmaBag.audio.length > 0}>
          <div className="flex flex-row items-start justify-between">
            <Hiragana word={lemmaBag.slug} />
            {lemmaBag.audio.length > 0 && <AudioControls audios={lemmaBag.audio} />}
          </div>
        </LemmaCardItem>
        <LemmaCardItem render={lemmaBag.lemmas.length > 0} className="py-3">
          <Senses projectedLemmas={lemmaBag.lemmas} previewLemma={previewLemma} />
        </LemmaCardItem>
        <LemmaCardItem render={lemmaBag.forms.length > 0}>
          <LemmaForms forms={lemmaBag.forms} />
        </LemmaCardItem>
        <LemmaCardItem render={lemmaBag.tags.length > 0}>
          <Tags tags={lemmaBag.tags} />
        </LemmaCardItem>
      </div>
    </div>
  );
}

function LemmaCardItem({
  render,
  className = 'py-7',
  children,
}: {
  render: boolean;
  className?: string;
  children: ReactNode;
}) {
  if (!render) {
    return null;
  }
  return <div className={className}>{children}</div>;
}

function LemmaTitle(props: { word: Word; toast: ToastFunction }) {
  const { word, toast } = props;
  const copyKanji = () => {
    navigator.clipboard.writeText(word.word);
    toast('Copied ' + word.word);
  };
  return (
    <div>
      <button
        onClick={(e) => {
          e.stopPropagation();
          copyKanji();
        }}
        className="transition-colors duration-300 hover:text-blue active:text-dark-blue"
      >
        <h1 className="text-6xl">
          <RenderWord word={word} />
        </h1>
      </button>
    </div>
  );
}

function Hiragana({ word }: { word: Word }) {
  return (
    <div className="text-2xl leading-10">
      {word.pitchShapes.length == 0
        ? word.furigana.map((e) => e.hiragana).join('')
        : word.pitchShapes.map((pitch, index) => {
            return (
              <span
                key={index}
                className={clsx(
                  'border-blue',
                  pitch.directions.includes(AccentDirection.Up) && 'border-t',
                  pitch.directions.includes(AccentDirection.Down) && 'border-b',
                  pitch.directions.includes(AccentDirection.Right) && 'border-r',
                  pitch.directions.includes(AccentDirection.Left) && 'border-l',
                )}
              >
                {pitch.hiragana}
              </span>
            );
          })}
    </div>
  );
}

function Senses({
  projectedLemmas,
  previewLemma,
}: {
  projectedLemmas: LemmaNoteInfo[];
  previewLemma: (lemma: LemmaNoteInfo) => Promise<void>;
}) {
  return (
    <ol className="-mx-4 text-lg">
      {projectedLemmas.map((projectedLemma, index) => {
        const lemma = projectedLemma.lemma;
        return (
          <li
            className="mb-3 flex flex-row justify-between rounded-xl p-4 last:mb-0 hover:bg-[#fdfdfd]"
            key={index}
          >
            <div>
              <p className="text-blue">{lemma.partsOfSpeech.join(', ')}</p>
              <p className="text-2xl">
                {index + 1}. {lemma.definitions.join('; ')}
              </p>
              <p className="text-dark-gray">{lemma.senseTags.join(', ')}</p>
            </div>
            <div
              className={clsx(
                'flex w-max shrink-0 flex-col place-content-center gap-1 pl-16',
                projectedLemma.noteID ? 'text-blue/60' : 'text-blue',
              )}
            >
              <button
                className="transition-color underline underline-offset-4 transition-colors duration-300 hover:text-green active:text-dark-green"
                onClick={() => previewLemma(projectedLemma)}
              >
                Add to Anki
              </button>
              <div className="self-center text-xs">
                {projectedLemma.noteID && '(already exists)'}
              </div>
            </div>
          </li>
        );
      })}
    </ol>
  );
}

function LemmaForms({ forms }: { forms: Word[] }) {
  return (
    <>
      <p className="text-lg text-blue">Other Forms</p>
      <div className="flex flex-col text-2xl">
        {forms.map((form, index) => {
          return (
            <p key={index}>
              {form.word}
              {form.hiragana ? '「' + form.hiragana + '」' : ''}
            </p>
          );
        })}
      </div>
    </>
  );
}

function Tags({ tags }: { tags: string[] }) {
  return (
    <>
      <p className="text-lg text-blue">Tags</p>
      <div className="flex flex-row justify-items-center gap-7 text-lg">
        {tags.map((tag, index) => (
          <p key={index} className="text-center">
            {tag}
          </p>
        ))}
      </div>
    </>
  );
}

function RenderFurigana(props: { furigana: Furigana }) {
  const { furigana } = props;
  if (furigana.kanji) {
    return (
      <>
        {furigana.kanji}
        <rp>[</rp>
        <rt className="text-lg text-dark-gray">{furigana.hiragana}</rt>
        <rp>]</rp>
      </>
    );
  } else {
    return (
      <>
        {furigana.hiragana}
        <rp>[</rp>
        <rt />
        <rp>]</rp>
      </>
    );
  }
}

function RenderWord(props: { word: Word }) {
  const { word } = props;
  if (word.furigana.length == 0) {
    return <>{word.word}</>;
  } else {
    return (
      <ruby>
        {word.furigana.map((furigana, index) => (
          // use index as key, because word.furigana will never be mutated
          <RenderFurigana key={index} furigana={furigana} />
        ))}
      </ruby>
    );
  }
}

function AudioControls({ audios }: { audios: Audio[] }) {
  const audio = useRef<HTMLAudioElement>(null);
  const handleMouseUp = () => {
    if (!audio.current) {
      return;
    }
    audio.current.play();
  };
  return (
    <div className="flex flex-row gap-2">
      <button
        onClick={(e) => {
          e.stopPropagation();
          handleMouseUp();
        }}
      >
        <PlayIcon className="transition-color h-10 w-10 fill-blue transition-colors duration-300 hover:fill-green active:fill-dark-green" />
      </button>
      <p className="my-auto text-lg leading-5 text-dark-gray">
        Listen to
        <br />
        pronunciation
      </p>
      {/* eslint-disable-next-line jsx-a11y/media-has-caption
       */}
      <audio ref={audio}>
        {audios.map((audio, index) => {
          return <source key={index} src={audio.source} type={audio.mediaType} />;
        })}
      </audio>
    </div>
  );
}
