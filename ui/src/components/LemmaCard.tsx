import { Fragment, useRef } from 'react';
import { Furigana, Lemma, Word, PitchType, Sense, Audio } from '../api/__generated__/graphql';
import { Id, toast } from 'react-toastify';
import { clsx } from 'clsx';

function LemmaTitle(props: { word: Word }) {
  const { word } = props;
  const render_furigana = (furigana: Furigana) => {
    if (furigana.kanji) {
      return (
        <>
          {furigana.kanji}
          <rp>[</rp>
          <rt className="text-dark-gray text-xl">{furigana.hiragana}</rt>
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
  };
  const render_word = () => {
    if (word.furigana.length == 0) {
      return <>{word.word}</>;
    } else {
      return (
        <ruby>
          {word.furigana.map((furigana, index) => (
            // use index as key, because word.furigana will never be mutated
            <Fragment key={index}>{render_furigana(furigana)}</Fragment>
          ))}
        </ruby>
      );
    }
  };
  // keep toast_id to show only one notification related to kanji copying
  // TODO: probably should get it from LemmaList or keep in context, because
  // different lemma cards can show notifications with different toast_id
  // defining toast_id as static string will not work, because seems that
  // it can not be dismissed and created at the same time
  const toast_id = useRef<Id | null>(null);
  const copy_kanji = () => {
    navigator.clipboard.writeText(word.word);
    if (toast_id.current) {
      toast.dismiss(toast_id.current);
    }
    toast_id.current = toast('Copied ' + word.word, {
      position: 'bottom-right',
      autoClose: 2000,
      type: 'info',
    });
  };
  return (
    <button
      onMouseUp={copy_kanji}
      className="hover:text-blue active:text-dark-blue transition-colors duration-300"
    >
      <h1 className="text-5xl pb-4">{render_word()}</h1>
    </button>
  );
}

function Hiragana(props: { word: Word }) {
  const { word } = props;
  if (word.hiragana == '') {
    return;
  }

  return (
    <div className="text-xl py-4">
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
    <div className="text-xl py-4">
      {senses.map((sense, index) => {
        return (
          <div className="mb-4 last:mb-0" key={index}>
            <p className="text-base text-blue">{sense.partOfSpeech.join(', ')}</p>
            <p>
              {index + 1}. {sense.definition.join('; ')}
            </p>
            <p className="text-base text-dark-gray">{sense.tags.join(', ')}</p>
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
    <div className="text-xl pt-4">
      <p className="text-base text-blue">Other Forms</p>
      <div className="flex flex-col">
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
    <div className="flex flex-col justify-items-center gap-4 py-4 text-lg text-blue">
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
  if (audios.length == 0) {
    return;
  }
  const audio = useRef<HTMLAudioElement>(null);
  const play_mouse_up = () => {
    if (!audio.current) {
      return;
    }
    audio.current.play();
  };
  return (
    <div className="flex flex-row pt-4">
      <button className="m-2" onMouseUp={play_mouse_up}>
        <svg
          className="fill-blue hover:fill-green active:fill-dark-green transition-color transition-colors duration-300"
          width="40"
          height="40"
          viewBox="10 10 41 41"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <circle cx="31" cy="31" r="20" />
          <path
            d="M26.5904 21.7116C26.7411 21.6368 26.9118 21.6579 27.0393 21.7559L41.0688 30.7172C41.1679 30.7994 41.2311 30.9185 41.2311 31.0471C41.2311 31.1735 41.1669 31.2968 41.0688 31.378L27.0383 40.3382C26.9645 40.3983 26.8675 40.4299 26.7748 40.4299L26.5893 40.3878C26.446 40.3193 26.3553 40.1696 26.3553 40.0073V22.09C26.3553 21.9298 26.446 21.7791 26.5904 21.7116Z"
            fill="white"
          />
        </svg>
      </button>
      <p className="m-1 text-dark-gray my-auto text-sm leading-5">Listen to pronunciation</p>
      <audio ref={audio}>
        {audios.map((audio, index) => {
          return <source key={index} src={audio.source} type={audio.type} />;
        })}
      </audio>
    </div>
  );
}

export default function LemmaCard(props: { lemma: Lemma }) {
  const { lemma } = props;

  return (
    <div className="flex flex-col sm:flex-row justify-between shadow-md rounded-md bg-gray my-4 px-8 pt-10 pb-7">
      <div className="grow divide-y divide-blue flex-1">
        <LemmaTitle word={lemma.slug} />
        <Hiragana word={lemma.slug} />
        <Senses senses={lemma.senses} />
        <LemmaForms forms={lemma.forms} />
      </div>
      <div className="shrink sm:basis-8 md:basis-14" />
      <div className="flex-none flex flex-col-reverse sm:flex-col justify-items-stretch basis-20 sm:basis-40">
        <button className="px-2 py-4 text-xl rounded-md bg-blue hover:bg-green active:bg-dark-green transition-colors duration-200 text-white">
          Add to Anki
        </button>
        <div className="flex flex-col-reverse sm:flex-col sm:divide-y divide-blue">
          <Tags tags={lemma.tags} />
          <AudioControls audios={lemma.audio} />
        </div>
      </div>
    </div>
  );
}
