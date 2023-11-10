import * as Dialog from '@radix-ui/react-dialog';
import { Cross1Icon } from '@radix-ui/react-icons';
import { clsx } from 'clsx';
import { ReactNode } from 'react';

export type DialogProps = {
  children: ReactNode;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  onCloseAutoFocus?: (event: Event) => void;
};

export function DialogModal({ children, open, onOpenChange, onCloseAutoFocus }: DialogProps) {
  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed h-full w-full top-0 left-0 overflow-y-auto overflow-x-hidden bg-blue opacity-80" />
        <Dialog.Content
          className={clsx(
            'fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-gray p-10',
            'w-[90%] sm:w-[calc(0.90*640px)]',
          )}
          onCloseAutoFocus={onCloseAutoFocus}
        >
          {children}
          <Dialog.Close asChild>
            <button className="absolute right-5 top-5" aria-label="Close">
              <Cross1Icon width="1.25rem" height="1.25rem" />
            </button>
          </Dialog.Close>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}