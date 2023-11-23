import { clsx } from 'clsx';
import { ForwardedRef, forwardRef } from 'react';

import { COLORS } from '../../colors';
import { IconProps } from './IconProps';

export const LoadingIcon = forwardRef<SVGSVGElement, IconProps>(function LoadinIcon(
  { colorInner = 'white', colorOuter = COLORS.blue, className, ...rest },
  forwardedRef,
) {
  return (
    <svg
      className={clsx('animate-rspin', className)}
      width="15"
      height="15"
      viewBox="0 0 45 45"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...rest}
      ref={forwardedRef}
    >
      <circle cx="22.5" cy="22.5" r="22.5" fill={colorOuter} />
      <path
        d="M9 16.8387C9 16.8387 12.1765 8.12903 22.5 8.12903C32.8235 8.12903 36 16.8387 36 16.8387"
        stroke={colorInner}
        strokeWidth="6"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path d="M9 16.8387H15.9677" stroke={colorInner} strokeWidth="6" strokeLinecap="round" />
      <path
        d="M9 28.1613C9 28.1613 12.1765 36.871 22.5 36.871C32.8235 36.871 36 28.1613 36 28.1613H29.0323"
        stroke={colorInner}
        strokeWidth="6"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
});

export const OKIcon = forwardRef<SVGSVGElement, IconProps>(function OKIcon(
  { colorInner = 'white', colorOuter = COLORS.green, ...rest },
  forwardedRef,
) {
  return (
    <svg
      viewBox="0 0 45 45"
      width="15"
      height="15"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...rest}
      ref={forwardedRef}
    >
      <path
        d="M37.0901 12.0787C38.4964 13.3021 38.6447 15.4339 37.4214 16.8402L22.7415 33.7152C22.1489 34.3963 21.3082 34.8116 20.4072 34.8684C19.5062 34.9251 18.6201 34.6185 17.9468 34.0171L6.75165 24.0171C5.36152 22.7754 5.24122 20.6418 6.48295 19.2517C7.72468 17.8616 9.85822 17.7413 11.2484 18.983L19.8932 26.7049L32.3287 12.4099C33.552 11.0036 35.6838 10.8553 37.0901 12.0787Z"
        fill={colorInner}
      />
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M22.5 45C34.9264 45 45 34.9264 45 22.5C45 10.0736 34.9264 0 22.5 0C10.0736 0 0 10.0736 0 22.5C0 34.9264 10.0736 45 22.5 45ZM37.4213 16.8401C38.6447 15.4338 38.4964 13.302 37.0901 12.0787C35.6838 10.8553 33.552 11.0036 32.3287 12.4099L19.8932 26.7049L11.2483 18.9829C9.85822 17.7412 7.72468 17.8615 6.48295 19.2517C5.24122 20.6418 5.36152 22.7753 6.75165 24.0171L17.9468 34.0171C18.6201 34.6185 19.5062 34.9251 20.4072 34.8683C21.3082 34.8116 22.1489 34.3962 22.7415 33.7151L37.4213 16.8401Z"
        fill={colorOuter}
      />
    </svg>
  );
});

export const WarningIcon = forwardRef<SVGSVGElement, IconProps>(function WarningIcon(
  { colorInner = 'white', colorOuter = COLORS.warningYellow, ...rest },
  forwardedRef,
) {
  return (
    <svg
      viewBox="0 0 45 45"
      width="15"
      height="15"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...rest}
      ref={forwardedRef}
    >
      <path
        d="M12.9914 21.4657C11.5489 21.7848 10.3375 22.9715 9.71263 24.1166C8.81981 25.7528 6.76962 26.3555 5.1334 25.4627C3.49717 24.5699 2.89452 22.5197 3.78733 20.8835C5.01544 18.6327 7.60258 15.7445 11.5336 14.875C15.7055 13.9523 20.2947 15.5235 24.886 20.1131C28.4476 23.6734 30.807 23.8001 32.0086 23.5343C33.4511 23.2153 34.6625 22.0286 35.2873 20.8835C36.1801 19.2472 38.2303 18.6446 39.8666 19.5374C41.5028 20.4302 42.1054 22.4804 41.2126 24.1166C39.9845 26.3673 37.3974 29.2556 33.4664 30.125C29.2944 31.0478 24.7053 29.4766 20.1139 24.887C16.5523 21.3267 14.193 21.2 12.9914 21.4657Z"
        fill={colorInner}
      />
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M22.5 45C34.9264 45 45 34.9264 45 22.5C45 10.0736 34.9264 0 22.5 0C10.0736 0 0 10.0736 0 22.5C0 34.9264 10.0736 45 22.5 45ZM9.71263 24.1166C10.3375 22.9715 11.5489 21.7848 12.9914 21.4657C14.193 21.1999 16.5523 21.3266 20.1139 24.8869C24.7053 29.4766 29.2944 31.0478 33.4664 30.125C37.3974 29.2555 39.9845 26.3673 41.2126 24.1166C42.1054 22.4804 41.5028 20.4302 39.8666 19.5374C38.2303 18.6445 36.1801 19.2472 35.2873 20.8834C34.6625 22.0285 33.4511 23.2152 32.0086 23.5343C30.807 23.8001 28.4476 23.6734 24.886 20.1131C20.2947 15.5234 15.7055 13.9522 11.5336 14.875C7.60257 15.7445 5.01544 18.6327 3.78733 20.8834C2.89452 22.5196 3.49717 24.5698 5.13339 25.4626C6.76962 26.3555 8.81981 25.7528 9.71263 24.1166Z"
        fill={colorOuter}
      />
    </svg>
  );
});

export const ErrorIcon = forwardRef<SVGSVGElement, IconProps>(function ErrorIcon(
  { colorInner = 'white', colorOuter = COLORS.lightRed, ...rest },
  forwardedRef,
) {
  return (
    <svg
      viewBox="0 0 45 45"
      width="15"
      height="15"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...rest}
      ref={forwardedRef}
    >
      <circle cx="22.5" cy="22.5" r="22.5" fill={colorOuter} />
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M22.5 31.25C26.25 31.25 27 11.9853 27 9.5C27 7.01472 24.9853 5 22.5 5C20.0147 5 18 7.01472 18 9.5C18 11.9853 18.75 31.25 22.5 31.25ZM22.5 41C24.5711 41 26.25 39.3211 26.25 37.25C26.25 35.1789 24.5711 33.5 22.5 33.5C20.4289 33.5 18.75 35.1789 18.75 37.25C18.75 39.3211 20.4289 41 22.5 41Z"
        fill={colorInner}
      />
    </svg>
  );
});

function getIconComponent(
  kind: StatusIconKind,
  props: IconProps,
  forwardedRef: ForwardedRef<SVGSVGElement>,
) {
  switch (kind) {
    case 'Loading':
      return <LoadingIcon {...props} ref={forwardedRef} />;
    case 'OK':
      return <OKIcon {...props} ref={forwardedRef} />;
    case 'Warning':
      return <WarningIcon {...props} ref={forwardedRef} />;
    case 'Error':
      return <ErrorIcon {...props} ref={forwardedRef} />;
  }
}

export type StatusIconKind = 'Loading' | 'OK' | 'Error' | 'Warning';

export type StatusIconProps = {
  kind: StatusIconKind;
} & IconProps;

export const StatusIcon = forwardRef<SVGSVGElement, StatusIconProps>(function StatusIcon(
  { kind, ...rest },
  forwardedRef,
) {
  const iconComponent = getIconComponent(kind, { ...rest }, forwardedRef);
  return iconComponent;
});

export default StatusIcon;
