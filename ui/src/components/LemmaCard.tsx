import { clsx } from 'clsx';
import { ReactNode, useRef } from 'react';

import { AccentDirection, Audio, Furigana, Lemma, Sense, Word } from '../api/__generated__/graphql';
import { ToastFunction } from '../lib/styled-toast';
import PlayIcon from './Icons/PlayIcon';

export type LemmaCardProps = {
  lemma: Lemma;
  toast: ToastFunction;
  previewLemma: (lemma: Lemma, senseIndex: number) => Promise<void>;
};

export default function LemmaCard({ lemma, toast, previewLemma }: LemmaCardProps) {
  const preview = (senseIndex: number) => previewLemma(lemma, senseIndex);
  return (
    <div className="shadow-md rounded-md bg-gray my-4 px-10 py-2">
      <div className="flex flex-col divide-y divide-blue">
        <LemmaCardItem render={true}>
          <LemmaTitle word={lemma.slug} toast={toast} />
        </LemmaCardItem>
        <LemmaCardItem render={lemma.slug.hiragana !== '' || lemma.audio.length > 0}>
          <div className="flex flex-row justify-between items-start">
            <Hiragana word={lemma.slug} />
            {lemma.audio.length > 0 && <AudioControls audios={lemma.audio} />}
          </div>
        </LemmaCardItem>
        <LemmaCardItem render={lemma.senses.length > 0} className="py-3">
          <Senses senses={lemma.senses} previewLemma={preview} />
        </LemmaCardItem>
        <LemmaCardItem render={lemma.forms.length > 0}>
          <LemmaForms forms={lemma.forms} />
        </LemmaCardItem>
        <LemmaCardItem render={lemma.tags.length > 0}>
          <Tags tags={lemma.tags} />
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
        className="hover:text-blue active:text-dark-blue transition-colors duration-300"
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
  senses,
  previewLemma,
}: {
  senses: Sense[];
  previewLemma: (senseIndex: number) => Promise<void>;
}) {
  return (
    <ol className="text-lg -mx-4">
      {senses.map((sense, index) => {
        return (
          <li
            className="mb-3 p-4 last:mb-0 hover:bg-[#fdfdfd] rounded-xl flex flex-row justify-between"
            key={index}
          >
            <div>
              <p className="text-blue">{sense.partOfSpeech.join(', ')}</p>
              <p className="text-2xl">
                {index + 1}. {sense.definition.join('; ')}
              </p>
              <p className="text-dark-gray">{sense.tags.join(', ')}</p>
            </div>
            <div className="pl-16 w-max shrink-0 flex items-center">
              <button
                className="h-full text-blue underline underline-offset-4 hover:text-green active:text-dark-green transition-color transition-colors duration-300"
                onClick={() => previewLemma(index)}
              >
                Add to Anki
              </button>
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
      <div className="text-2xl flex flex-col">
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
        <rt className="text-dark-gray text-lg">{furigana.hiragana}</rt>
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
        <PlayIcon className="h-10 w-10 fill-blue hover:fill-green active:fill-dark-green transition-color transition-colors duration-300" />
      </button>
      <p className="text-dark-gray my-auto text-lg leading-5">
        Listen to
        <br />
        pronunciation
      </p>
      {/* eslint-disable-next-line jsx-a11y/media-has-caption
       */}
      <audio ref={audio}>
        {audios.map((audio, index) => {
          return <source key={index} src={audio.source} type={audio.type} />;
        })}
      </audio>
    </div>
  );
}
