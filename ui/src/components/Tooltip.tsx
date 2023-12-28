import * as RadixTooltip from '@radix-ui/react-tooltip';
import { clsx } from 'clsx';
import { ReactNode } from 'react';

type TooltipProps = {
  children: ReactNode;
  content: ReactNode;
};

export default function Tooltip({ children, content }: TooltipProps) {
  return (
    <RadixTooltip.Provider>
      <RadixTooltip.Root>
        <RadixTooltip.Trigger asChild>{children}</RadixTooltip.Trigger>
        <RadixTooltip.Content
          className={clsx(
            `text-md select-none rounded bg-white px-3.5 py-2.5 leading-none text-dark-gray 
             shadow-[hsl(206_22%_7%_/_35%)_0px_10px_38px_-10px,_hsl(206_22%_7%_/_20%)_0px_10px_20px_-15px]
             will-change-[transform,opacity]
             radix-state-delayed-open:radix-side-bottom:animate-slideUpAndFade
             radix-state-delayed-open:radix-side-left:animate-slideRightAndFade
             radix-state-delayed-open:radix-side-right:animate-slideLeftAndFade
             radix-state-delayed-open:radix-side-top:animate-slideDownAndFade`,
          )}
          sideOffset={0}
        >
          {content}
          <RadixTooltip.Arrow className="fill-white" />
        </RadixTooltip.Content>
      </RadixTooltip.Root>
    </RadixTooltip.Provider>
  );
}
