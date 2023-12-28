import * as Dialog from '@radix-ui/react-dialog';
import { Cross1Icon } from '@radix-ui/react-icons';
import { clsx } from 'clsx';
import { ReactNode } from 'react';

export enum DialogWidth {
  Small = 'dialog-sm',
  Medium = 'dialog-md',
  Large = 'dialog-lg',
}

export type DialogProps = {
  children: ReactNode;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  onCloseAutoFocus?: (event: Event) => void;

  widthVariant?: DialogWidth;
};

export function DialogModal({
  children,
  open,
  onOpenChange,
  onCloseAutoFocus,
  widthVariant,
}: DialogProps) {
  widthVariant = widthVariant ?? DialogWidth.Small;

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed bottom-0 left-0 right-0 top-0 grid place-items-center overflow-y-auto overflow-x-hidden bg-blue/60 backdrop-blur-sm">
          <Dialog.Content
            className={clsx('relative bg-gray p-10', widthVariant)}
            onCloseAutoFocus={onCloseAutoFocus}
          >
            <Dialog.Close asChild>
              <button className="absolute right-5 top-5" aria-label="Close">
                <Cross1Icon className="hover:drop-shadow" width="1.25rem" height="1.25rem" />
              </button>
            </Dialog.Close>
            {children}
          </Dialog.Content>
        </Dialog.Overlay>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
