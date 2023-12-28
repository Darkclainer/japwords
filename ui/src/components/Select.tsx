import { CheckIcon, ChevronDownIcon, ChevronUpIcon, TriangleDownIcon } from '@radix-ui/react-icons';
import type * as Radix from '@radix-ui/react-primitive';
import * as Select from '@radix-ui/react-select';
import clsx from 'clsx';
import { forwardRef } from 'react';

import { COLORS } from '../colors';

export type SelectTriggerProps = Radix.ComponentPropsWithoutRef<typeof Select.Trigger> & {
  hasError?: boolean;
};

export const SelectTrigger = forwardRef<
  React.ElementRef<typeof Select.Trigger>,
  SelectTriggerProps
>(function SelectTrigger({ children, className, hasError, ...props }, forwardedRef) {
  return (
    <Select.Trigger
      data-error={hasError || null}
      className={clsx(
        `group/trigger flex flex-row justify-between overflow-x-hidden bg-white
         px-2 py-1 text-xl text-black ring-1 ring-blue
         data-error:text-error-red data-error:ring-error-red data-disabled:text-dark-gray data-disabled:ring-dark-gray`,
        className,
      )}
      {...props}
      ref={forwardedRef}
    >
      {children}
      <Select.Icon className="">
        <TriangleDownIcon className="inline text-black group-data-error/trigger:text-error-red group-data-disabled/trigger:text-dark-gray group-data-state-open/trigger:rotate-180" />
      </Select.Icon>
    </Select.Trigger>
  );
});

export const SelectContent = forwardRef<
  React.ElementRef<typeof Select.Content>,
  Radix.ComponentPropsWithoutRef<typeof Select.Content>
>(function SelectContent({ children, className, ...props }, forwardedRef) {
  return (
    <Select.Content
      className={clsx(
        `group/content bg-white
         text-xl ring-1 ring-blue
         max-h-radix-select-content-available-height
         w-radix-select-trigger-width`,
        className,
      )}
      sideOffset={0}
      position="popper"
      {...props}
      ref={forwardedRef}
    >
      {children}
    </Select.Content>
  );
});

export const SelectSeparator = forwardRef<
  React.ElementRef<typeof Select.Separator>,
  Radix.ComponentPropsWithoutRef<typeof Select.Separator>
>(function SelectSeparator({ children, className, ...props }, forwardedRef) {
  return (
    <Select.Separator
      className={clsx('mx-2 h-px bg-blue', className)}
      {...props}
      ref={forwardedRef}
    >
      {children}
    </Select.Separator>
  );
});

export const SelectScrollUpButton = forwardRef<
  React.ElementRef<typeof Select.ScrollUpButton>,
  Radix.ComponentPropsWithoutRef<typeof Select.ScrollUpButton>
>(function SelectScrollUpButton({ children, className, ...props }, forwardedRef) {
  return (
    <Select.ScrollUpButton
      className={clsx('flex items-center justify-center', className)}
      {...props}
      ref={forwardedRef}
    >
      {children ?? <ChevronUpIcon color={COLORS.blue} />}
    </Select.ScrollUpButton>
  );
});

export const SelectScrollDownButton = forwardRef<
  React.ElementRef<typeof Select.ScrollDownButton>,
  Radix.ComponentPropsWithoutRef<typeof Select.ScrollDownButton>
>(function SelectScrollDownButton({ children, className, ...props }, forwardedRef) {
  return (
    <Select.ScrollDownButton
      className={clsx('flex items-center justify-center', className)}
      {...props}
      ref={forwardedRef}
    >
      {children ?? <ChevronDownIcon color={COLORS.blue} />}
    </Select.ScrollDownButton>
  );
});

export const SelectItem = forwardRef<
  React.ElementRef<typeof Select.Item>,
  Radix.ComponentPropsWithoutRef<typeof Select.Item>
>(function SelectItem({ children, className, ...props }, forwardedRef) {
  return (
    <Select.Item
      className={clsx(
        `oveflow-x-hidden m-1 flex flex-row
         items-center justify-between
         px-1
         radix-highlighted:text-blue`,
        className,
      )}
      {...props}
      ref={forwardedRef}
    >
      {children}
      <Select.ItemIndicator>
        <CheckIcon color={COLORS.blue} />
      </Select.ItemIndicator>
    </Select.Item>
  );
});
