import clsx from 'clsx';
import { ComponentPropsWithoutRef, forwardRef } from 'react';

export type TextInputProps = ComponentPropsWithoutRef<'input'> & {
  hasError?: boolean;
};

const TextInput = forwardRef<React.ElementRef<'input'>, TextInputProps>(function TextInput(
  { hasError, className, ...other },
  forwardedRef,
) {
  return (
    <input
      data-error={hasError || null}
      className={clsx(
        'appearance-none py-1 px-2 shadow-sm',
        'leading-6 placeholder-slate-400',
        'text-xl',
        'focus:!text-black data-error:text-error-red',
        'ring-1 focus:ring-[3px] focus:outline-none',
        'ring-blue data-error:ring-error-red',
        className,
      )}
      {...other}
      ref={forwardedRef}
    />
  );
});

export default TextInput;
