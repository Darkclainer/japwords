import { clsx } from 'clsx';
import { useRef } from 'react';

import { Audio, Furigana, Lemma, PitchType, Sense, Word } from '../api/__generated__/graphql';
import { ToastFunction } from '../lib/styled-toast';
import Button from './Button';
import PlayIcon from './Icons/PlayIcon';

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

function LemmaTitle(props: { word: Word; toast: ToastFunction }) {
  const { word, toast } = props;
  const copyKanji = () => {
    navigator.clipboard.writeText(word.word);
    toast('Copied ' + word.word);
  };
  return (
    <button
      onClick={(e) => {
        e.stopPropagation();
        copyKanji();
      }}
      className="hover:text-blue active:text-dark-blue transition-colors duration-300"
    >
      <h1 className="text-6xl pb-5">
        <RenderWord word={word} />
      </h1>
    </button>
  );
}

function Hiragana(props: { word: Word }) {
  const { word } = props;
  if (word.hiragana == '') {
    return;
  }

  return (
    <div className="text-2xl py-5">
      {word.pitch.length == 0
        ? word.furigana.map((e) => e.hiragana).join('')
        : word.pitch.map((pitch, index) => {
            return (
              <span
                key={index}
                className={clsx(
                  'border-blue',
                  pitch.pitch.includes(PitchType.Up) && 'border-t',
                  pitch.pitch.includes(PitchType.Down) && 'border-b',
                  pitch.pitch.includes(PitchType.Right) && 'border-r',
                  pitch.pitch.includes(PitchType.Left) && 'border-l',
                )}
              >
                {pitch.hiragana}
              </span>
            );
          })}
    </div>
  );
}

function Senses(props: { senses: Sense[] }) {
  const { senses } = props;
  if (senses.length == 0) {
    return;
  }
  return (
    <div className="text-lg py-7">
      {senses.map((sense, index) => {
        return (
          <div className="mb-7 last:mb-0" key={index}>
            <p className="text-blue">{sense.partOfSpeech.join(', ')}</p>
            <p className="text-2xl">
              {index + 1}. {sense.definition.join('; ')}
            </p>
            <p className="text-dark-gray">{sense.tags.join(', ')}</p>
          </div>
        );
      })}
    </div>
  );
}

function LemmaForms(props: { forms: Word[] }) {
  const { forms } = props;
  if (forms.length == 0) {
    return;
  }
  return (
    <div className="pt-7">
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
    </div>
  );
}

function Tags(props: { tags: string[] }) {
  const { tags } = props;
  if (tags.length == 0) {
    return;
  }
  return (
    <div className="flex flex-col justify-items-center gap-7 py-5 text-2xl text-blue">
      {tags.map((tag, index) => {
        return (
          <p key={index} className="text-center">
            {tag}
          </p>
        );
      })}
    </div>
  );
}

function AudioControls(props: { audios: Audio[] }) {
  const { audios } = props;
  const audio = useRef<HTMLAudioElement>(null);
  const handleMouseUp = () => {
    if (!audio.current) {
      return;
    }
    audio.current.play();
  };
  if (audios.length == 0) {
    return;
  }
  return (
    <div className="flex flex-row pt-5">
      <button
        className="m-2"
        onClick={(e) => {
          e.stopPropagation();
          handleMouseUp();
        }}
      >
        <PlayIcon className="h-10 w-10 fill-blue hover:fill-green active:fill-dark-green transition-color transition-colors duration-300" />
      </button>
      <p className="m-1 text-dark-gray my-auto text-lg leading-5">Listen to pronunciation</p>
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

export default function LemmaCard(props: { lemma: Lemma; toast: ToastFunction }) {
  const { lemma, toast } = props;

  return (
    <div className="flex flex-col sm:flex-row justify-between shadow-md rounded-md bg-gray my-4 px-10 pt-12 pb-7">
      <div className="grow divide-y divide-blue flex-1">
        <LemmaTitle word={lemma.slug} toast={toast} />
        <Hiragana word={lemma.slug} />
        <Senses senses={lemma.senses} />
        <LemmaForms forms={lemma.forms} />
      </div>
      <div className="shrink sm:basis-8 md:basis-24" />
      <div className="flex-none flex flex-col-reverse sm:flex-col justify-items-stretch basis-20 sm:basis-52">
        <Button>Add to Anki</Button>
        <div className="flex flex-col-reverse sm:flex-col sm:divide-y divide-blue">
          <Tags tags={lemma.tags} />
          <AudioControls audios={lemma.audio} />
        </div>
      </div>
    </div>
  );
}
