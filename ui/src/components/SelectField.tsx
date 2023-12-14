import * as RSelect from '@radix-ui/react-select';
import { FieldHookConfig, useField } from 'formik';
import { ReactNode } from 'react';

import {
  SelectContent,
  SelectScrollDownButton,
  SelectScrollUpButton,
  SelectTrigger,
} from './Select';

export type SelectFieldProps = FieldHookConfig<string> & {
  id?: string;
  placeholderLabel: string;
  disabled?: boolean;
  onValueChange?: (value: string) => void;
  children: ReactNode;
};

export default function SelectField({
  children,
  name,
  id,
  disabled,
  placeholderLabel,
  onValueChange,
}: SelectFieldProps) {
  const [field, meta, helpers] = useField<string>(name);
  return (
    <RSelect.Root
      disabled={disabled}
      value={field.value}
      onValueChange={(value) => {
        if (onValueChange !== undefined) {
          onValueChange(value);
        }
        helpers.setValue(value);
      }}
    >
      <SelectTrigger id={id} className="w-full" hasError={!!meta.error}>
        <RSelect.Value placeholder={placeholderLabel} />
      </SelectTrigger>
      <RSelect.Portal>
        <SelectContent>
          <SelectScrollUpButton />
          <RSelect.Viewport>{children}</RSelect.Viewport>
          <SelectScrollDownButton />
        </SelectContent>
      </RSelect.Portal>
    </RSelect.Root>
  );
}
