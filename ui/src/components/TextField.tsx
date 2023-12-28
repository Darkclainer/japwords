import { FieldConfig, useField } from 'formik';
import TextInput from './TextInput';
import { ComponentPropsWithoutRef, forwardRef } from 'react';

export type TextFieldProps = FieldConfig<string> &
  ComponentPropsWithoutRef<'input'> & {
    inputClassName?: string;
  };

const TextField = forwardRef<React.ElementRef<'input'>, TextFieldProps>(
  function TextField(props, forwadedRef) {
    const [field, meta] = useField(props);
    const { inputClassName, ...rest } = props;
    return (
      <TextInput
        className={inputClassName}
        hasError={!!meta.error}
        {...field}
        {...rest}
        ref={forwadedRef}
      />
    );
  },
);

export default TextField;
