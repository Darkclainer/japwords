import { clsx } from 'clsx';
import { ComponentPropsWithoutRef } from 'react';

export enum ButtonVariant {
  Primary = 'PRIMARY',
  Dangerous = 'DANGEROUS',
}

interface ButtonProps extends ComponentPropsWithoutRef<'button'> {
  variant?: ButtonVariant;
}

export default function Button(props: ButtonProps) {
  const { variant = ButtonVariant.Primary, className, ...rest } = props;

  let variantClasses = '';
  switch (variant) {
    case ButtonVariant.Primary:
      variantClasses = 'bg-blue hover:bg-green active:bg-dark-green disabled:bg-dark-gray';
      break;
    case ButtonVariant.Dangerous:
      variantClasses = 'bg-red hover:bg-light-red active:bg-dark-red disabled:bg-dark-gray';
      break;

    default:
      throw new Error('unreachable');
  }

  return (
    <button
      className={clsx(
        'px-3 py-6 text-2xl',
        variantClasses,
        'text-white rounded-md transition-colors duration-200',
        className,
      )}
      {...rest}
    >
      {props.children}
    </button>
  );
}
