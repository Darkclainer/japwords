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
        'py-1 px-2',
        'bg-white',
        'ring-1 ring-blue data-disabled:ring-dark-gray data-error:ring-error-red',
        'text-xl text-black data-disabled:text-dark-gray data-error:text-error-red',
        'flex flex-row justify-between',
        'overflow-x-hidden',
        'group/trigger',
        className,
      )}
      {...props}
      ref={forwardedRef}
    >
      {children}
      <Select.Icon className="">
        <TriangleDownIcon
          className={clsx(
            'inline group-data-state-open/trigger:rotate-180 text-black group-data-disabled/trigger:text-dark-gray group-data-error/trigger:text-error-red',
          )}
        />
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
        'bg-white',
        'ring-1 ring-blue',
        'text-xl',
        'w-radix-select-trigger-width',
        'max-h-radix-select-content-available-height',
        'group/content',
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
      className={clsx('h-px bg-blue mx-2', className)}
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
        'flex flex-row items-center justify-between',
        'radix-highlighted:text-blue',
        'm-1',
        'px-1',
        'oveflow-x-hidden',
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
