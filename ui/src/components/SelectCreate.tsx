import * as Dialog from '@radix-ui/react-dialog';
import {
  CheckIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  PlusIcon,
  TriangleDownIcon,
} from '@radix-ui/react-icons';
import * as Label from '@radix-ui/react-label';
import type * as Radix from '@radix-ui/react-primitive';
import * as Select from '@radix-ui/react-select';
import { clsx } from 'clsx';
import { ErrorMessage, Form, Formik, FormikHelpers } from 'formik';
import React, { forwardRef, useCallback, useId, useMemo, useRef, useState } from 'react';
import { Result } from 'true-myth';
import * as yup from 'yup';

import { COLORS } from '../colors';
import Button from './Button';
import { DialogModal } from './DialogModal';
import TextField from './TextField';

export type SelectItem = {
  value: string;
  label?: string;
};

export type SelectCreateProps = {
  id: string;

  items: SelectItem[];
  selectedValue?: string;
  onValueChange: (value: string) => Promise<void>;
  handleCreate: (value: string) => Promise<Result<string, string>>;
  validateValue?: (value: string) => string | null;

  isDisabled?: boolean;
  hasError?: boolean;

  triggerClassName?: string;

  createLabel: string;
  createDefaultValue?: string;
  placeholderLabel: string;
  dialogTitle: string;
  dialogInputLabel: string;
};

type FormData = {
  name: string;
};

const createPlaceholder = '__create_placeholder';

enum State {
  Enabled,
  Disabled,
  Creating,
}

export default function SelectCreate({
  id,
  items,
  selectedValue,
  onValueChange,
  handleCreate,
  validateValue,
  isDisabled,
  hasError,
  triggerClassName,
  createLabel,
  createDefaultValue,
  placeholderLabel,
  dialogTitle,
  dialogInputLabel,
}: SelectCreateProps) {
  const [state, setState] = useState(State.Enabled);
  const actualState = isDisabled ? State.Disabled : state;
  const stateColor = (element: string, defaultColor?: string) =>
    actualState != State.Enabled
      ? `${element}-dark-gray`
      : hasError
      ? `${element}-error-red`
      : `${element}-${defaultColor ?? 'blue'}`;

  const triggerRef = useRef<HTMLButtonElement>(null);
  const validationSchema = useMemo(
    () =>
      yup.object({
        name: yup
          .string()
          .required()
          .test((value, context) => {
            if (value == createPlaceholder) {
              return context.createError({
                message: `name can not be reserved value: ${createPlaceholder}`,
              });
            }
            if (items.find((item) => item.value == value)) {
              return context.createError({
                message: 'name already exists',
              });
            }
            if (validateValue) {
              const error = validateValue(value);
              if (error) {
                return context.createError({
                  message: error,
                });
              }
            }
            return true;
          }),
      }),
    [items, validateValue],
  );
  const onValueChangeImpl = useCallback(
    (value: string) => {
      if (value == createPlaceholder) {
        setState(State.Creating);
      } else {
        setState(State.Disabled);
        onValueChange(value).finally(() => setState(State.Enabled));
      }
    },
    [onValueChange],
  );
  const handleSubmit = useCallback(
    async (values: FormData, helpers: FormikHelpers<FormData>) => {
      const result = await handleCreate(values.name);
      if (result.isOk) {
        await onValueChange(result.value);
        setState(State.Enabled);
      } else {
        helpers.setErrors({
          name: result.error,
        });
      }
    },
    [handleCreate, onValueChange],
  );
  const dialogInputId = useId();

  return (
    <>
      <Select.Root
        disabled={actualState !== State.Enabled}
        value={selectedValue || ''}
        onValueChange={onValueChangeImpl}
      >
        <Select.Trigger
          id={id}
          className={clsx(
            'py-1 px-2',
            'bg-white ring-1',
            stateColor('ring'),
            'text-xl',
            stateColor('text', 'black'),
            'group',
            'flex flex-row justify-between',
            'overflow-x-hidden',
            triggerClassName,
          )}
          ref={triggerRef}
        >
          <Select.Value placeholder={placeholderLabel} asChild>
            <span>{selectedValue}</span>
          </Select.Value>
          <Select.Icon className="">
            <TriangleDownIcon
              className={clsx('inline group-radix-state-open:rotate-180', stateColor('text'))}
            />
          </Select.Icon>
        </Select.Trigger>
        <Select.Portal>
          <Select.Content
            className={clsx(
              'bg-white',
              'ring-1 ring-blue',
              'text-xl',
              'w-radix-select-trigger-width',
              'max-h-radix-select-content-available-height',
              'group',
            )}
            sideOffset={5}
            position="popper"
          >
            <Select.Group className="group-radix-side-top:order-3">
              <SelectItem value={createPlaceholder}>
                <Select.ItemText>{createLabel}</Select.ItemText>
                <PlusIcon className="inline" color={COLORS.blue} />
              </SelectItem>
            </Select.Group>
            <Select.Separator className="h-px bg-blue mx-2 group-radix-side-top:order-2" />
            <Select.ScrollUpButton className="flex items-center justify-center">
              <ChevronUpIcon color={COLORS.blue} />
            </Select.ScrollUpButton>
            <Select.Viewport>
              <Select.Group className="order-1">
                {items.map((item) => (
                  <SelectItem key={item.value} value={item.value}>
                    <Select.ItemText>{item.label ?? item.value}</Select.ItemText>
                    <Select.ItemIndicator className="">
                      <CheckIcon color={COLORS.blue} />
                    </Select.ItemIndicator>
                  </SelectItem>
                ))}
              </Select.Group>
            </Select.Viewport>
            <Select.ScrollDownButton className="flex items-center justify-center">
              <ChevronDownIcon color={COLORS.blue} />
            </Select.ScrollDownButton>
          </Select.Content>
        </Select.Portal>
      </Select.Root>
      <DialogModal
        open={state === State.Creating}
        onOpenChange={(open) => setState(open ? State.Creating : State.Enabled)}
        onCloseAutoFocus={(event) => {
          triggerRef.current?.focus();
          event.preventDefault();
        }}
      >
        <Dialog.Title className="mb-2.5 text-2xl font-bold text-blue">{dialogTitle}</Dialog.Title>
        <Formik<FormData>
          initialValues={{ name: createDefaultValue ?? 'NewValue' }}
          onSubmit={handleSubmit}
          validationSchema={validationSchema}
        >
          {(props) => {
            return (
              <Form className="flex flex-col gap-2.5">
                <Label.Root className="text-blue text-lg" htmlFor={dialogInputId}>
                  {dialogInputLabel}
                </Label.Root>
                <TextField
                  name="name"
                  type="input"
                  inputClassName="my-2.5 w-full shrink-0"
                  id={dialogInputId}
                />
                <p className="text-error-red text-lg">
                  <ErrorMessage name="name" />
                </p>
                <Button type="submit" disabled={props.isSubmitting}>
                  Create
                </Button>
              </Form>
            );
          }}
        </Formik>
      </DialogModal>
    </>
  );
}

const SelectItem = forwardRef<
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
    </Select.Item>
  );
});
