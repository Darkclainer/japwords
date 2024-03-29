import * as Dialog from '@radix-ui/react-dialog';
import { PlusIcon } from '@radix-ui/react-icons';
import * as Label from '@radix-ui/react-label';
import * as RSelect from '@radix-ui/react-select';
import { ErrorMessage, Form, Formik, FormikHelpers } from 'formik';
import { useCallback, useId, useMemo, useRef, useState } from 'react';
import { Result } from 'true-myth';
import * as yup from 'yup';

import { COLORS } from '../colors';
import Button from './Button';
import { DialogModal } from './DialogModal';
import {
  SelectContent,
  SelectItem,
  SelectScrollDownButton,
  SelectScrollUpButton,
  SelectSeparator,
  SelectTrigger,
} from './Select';
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
      <RSelect.Root
        disabled={actualState !== State.Enabled}
        value={selectedValue || ''}
        onValueChange={onValueChangeImpl}
      >
        <SelectTrigger hasError={hasError} id={id} className={triggerClassName} ref={triggerRef}>
          <RSelect.Value placeholder={placeholderLabel} asChild>
            <span>{selectedValue}</span>
          </RSelect.Value>
        </SelectTrigger>
        <RSelect.Portal>
          <SelectContent>
            <RSelect.Group className="group-data-side-top/content:order-3">
              <SelectItem value={createPlaceholder}>
                <RSelect.ItemText>{createLabel}</RSelect.ItemText>
                <PlusIcon className="inline" color={COLORS.blue} />
              </SelectItem>
            </RSelect.Group>
            <SelectSeparator className="group-data-side-top/content:order-2" />
            <SelectScrollUpButton />
            <RSelect.Viewport>
              <RSelect.Group className="order-1">
                {items.map((item) => (
                  <SelectItem key={item.value} value={item.value}>
                    <RSelect.ItemText>{item.label ?? item.value}</RSelect.ItemText>
                  </SelectItem>
                ))}
              </RSelect.Group>
            </RSelect.Viewport>
            <SelectScrollDownButton />
          </SelectContent>
        </RSelect.Portal>
      </RSelect.Root>
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
                <Label.Root className="text-lg text-blue" htmlFor={dialogInputId}>
                  {dialogInputLabel}
                </Label.Root>
                <TextField
                  name="name"
                  type="input"
                  inputClassName="my-2.5 w-full shrink-0"
                  id={dialogInputId}
                />
                <p className="text-lg text-error-red">
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
