import { clsx } from 'clsx';
import { ReactNode, useEffect, useId, useRef, useState } from 'react';
import Editor from 'react-simple-code-editor';

type CodeEditorProps = {
  id?: string;
  value: string;
  onValueChange: (v: string) => void;
  onBlur?: () => void;
};

export default function CodeEditor({ id, value, onValueChange, onBlur }: CodeEditorProps) {
  const [inFocus, setInFocus] = useState(false);
  const generatedTextareaId = useId();
  const textareaId = id ?? generatedTextareaId;
  const textareaRef = useRef<HTMLTextAreaElement>();
  useEffect(() => {
    textareaRef.current = document.getElementById(textareaId) as HTMLTextAreaElement;
  });
  return (
    // As I understand, we don't need to add spurious events or aria labels because we handle onClick
    // to support mouse users, if they click on part of editor that's visually are not editor yet!
    // eslint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-static-element-interactions
    <div
      className={clsx(
        `overflow-y-scroll, max-h-60 min-h-[7.5rem]
         bg-white ring-blue focus:outline-none`,
        inFocus ? 'ring-[3px]' : 'ring-1',
      )}
      onClick={() => {
        // little fix for selection when we have on line and empty space
        if (inFocus || textareaRef.current === undefined) {
          return;
        }
        const end = textareaRef.current.value.length;
        textareaRef.current.setSelectionRange(end, end);
        textareaRef.current.focus();
      }}
    >
      <Editor
        className="min-h-full font-mono text-lg"
        preClassName="!pl-14 !break-all"
        textareaClassName="!pl-14 !break-all focus:outline-none"
        textareaId={textareaId}
        onFocus={() => setInFocus(true)}
        onBlur={() => {
          setInFocus(false);
          if (onBlur) {
            onBlur();
          }
        }}
        value={value}
        padding={10}
        onValueChange={onValueChange}
        highlight={(code) => hightlightWithNumbers(code, (v) => v)}
      />
    </div>
  );
}

function hightlightWithNumbers(input: string, hightlight: (src: string) => string): ReactNode {
  return (
    <>
      {hightlight(input)
        .split('\n')
        .map((line, i) => (
          <div key={i}>
            <span className="absolute left-0 h-6 w-10 text-right text-dark-gray">{i + 1}</span>
            {line}
            {'\n'}
          </div>
        ))}
    </>
  );
}
