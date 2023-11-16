import { clsx } from 'clsx';
import { ReactNode, useState } from 'react';
import Editor from 'react-simple-code-editor';

type CodeEditorProps = {
  value: string;
  onValueChange: (v: string) => void;
};

export default function CodeEditor({ value, onValueChange }: CodeEditorProps) {
  const [inFocus, setInFocus] = useState(false);
  return (
    <div
      className={clsx(
        'overflow-y-scroll min-h-[7.5rem] max-h-60',
        'bg-white ring-blue focus:outline-none',
        inFocus ? 'ring-[3px]' : 'ring-1',
      )}
    >
      <Editor
        className="text-lg min-h-full font-mono"
        preClassName="!pl-14"
        textareaClassName="!pl-14 focus:outline-none"
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
