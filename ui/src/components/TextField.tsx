import { clsx } from 'clsx';
import { FieldHookConfig, useField } from 'formik';

export type TextFieldProps = FieldHookConfig<string> &
  JSX.IntrinsicElements['input'] & {
    inputClassName?: string;
  };

export default function TextField(props: TextFieldProps) {
  const [field, meta] = useField(props);
  const { inputClassName, ...rest } = props;
  const isValid = !meta.error;
  return (
    <input
      className={clsx(
        'self-start',
        'ring-1 focus:ring-[3px] focus:outline-none',
        'appearance-none w-full py-1 px-2 shadow-sm',
        'text-xl',
        'focus:text-black',
        'leading-6 placeholder-slate-400',
        inputClassName,
        !isValid && 'text-error-red',
        isValid ? 'ring-blue' : 'ring-error-red',
      )}
      {...rest}
      {...field}
    />
  );
}
