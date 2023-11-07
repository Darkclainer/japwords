import { Label } from '@radix-ui/react-label';
import { useId } from 'react';

export function NoteSelect() {
  const noteTriggerId = useId();
  return (
    <div className="flex flex-col gap-5">
      <Label className="text-2xl" htmlFor={noteTriggerId}>
        Choose a note type:
      </Label>
    </div>
  );
}
