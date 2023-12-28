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
        `placeholder-slate-400 appearance-none px-2 py-1 text-xl leading-6 shadow-sm
         ring-1 ring-blue focus:!text-black focus:outline-none focus:ring-[3px]
         data-error:text-error-red data-error:ring-error-red`,
        className,
      )}
      {...other}
      ref={forwardedRef}
    />
  );
});

export default TextInput;
