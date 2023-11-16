import { clsx } from 'clsx';
import { ReactNode, useEffect, useId, useRef, useState } from 'react';
import Editor from 'react-simple-code-editor';

type CodeEditorProps = {
  value: string;
  onValueChange: (v: string) => void;
};

export default function CodeEditor({ value, onValueChange }: CodeEditorProps) {
  const [inFocus, setInFocus] = useState(false);
  const textareaId = useId();
  const textareaRef = useRef<HTMLTextAreaElement>();
  useEffect(() => {
    textareaRef.current = document.getElementById(textareaId) as HTMLTextAreaElement;
  });
  return (
    <div
      className={clsx(
        'overflow-y-scroll min-h-[7.5rem] max-h-60',
        'bg-white ring-blue focus:outline-none',
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
        className="text-lg min-h-full font-mono"
        preClassName="!pl-14 !break-all"
        textareaClassName="!pl-14 !break-all focus:outline-none"
        textareaId={textareaId}
        onFocus={() => setInFocus(true)}
        onBlur={() => setInFocus(false)}
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
            <span className="absolute text-dark-gray h-6 left-0 text-right w-10">{i + 1}</span>
            {line}
            {'\n'}
          </div>
        ))}
    </>
  );
}
