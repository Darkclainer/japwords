import { ComponentPropsWithoutRef } from 'react';

const defaultIconSize = 35;

export type StatusIconSizeProps = {
  size: number | string;
};

export function LoadingIcon({ size }: StatusIconSizeProps) {
  return (
    <svg
      className="animate-rspin"
      width={size}
      height={size}
      viewBox="0 0 45 45"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <circle cx="22.5" cy="22.5" r="22.5" fill="#617DB6" />
      <path
        d="M9 16.8387C9 16.8387 12.1765 8.12903 22.5 8.12903C32.8235 8.12903 36 16.8387 36 16.8387"
        stroke="white"
        strokeWidth="6"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path d="M9 16.8387H15.9677" stroke="white" strokeWidth="6" strokeLinecap="round" />
      <path
        d="M9 28.1613C9 28.1613 12.1765 36.871 22.5 36.871C32.8235 36.871 36 28.1613 36 28.1613H29.0323"
        stroke="white"
        strokeWidth="6"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}

export function OKIcon({ size }: StatusIconSizeProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 45 45"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M22.5 45C34.9264 45 45 34.9264 45 22.5C45 10.0736 34.9264 0 22.5 0C10.0736 0 0 10.0736 0 22.5C0 34.9264 10.0736 45 22.5 45ZM37.4213 16.8401C38.6447 15.4338 38.4964 13.302 37.0901 12.0787C35.6838 10.8553 33.552 11.0036 32.3287 12.4099L19.8932 26.7049L11.2483 18.9829C9.85822 17.7412 7.72468 17.8615 6.48295 19.2517C5.24122 20.6418 5.36152 22.7753 6.75165 24.0171L17.9468 34.0171C18.6201 34.6185 19.5062 34.9251 20.4072 34.8683C21.3082 34.8116 22.1489 34.3962 22.7415 33.7151L37.4213 16.8401Z"
        fill="#77B661"
      />
    </svg>
  );
}

export function WarningIcon({ size }: StatusIconSizeProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 45 45"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M22.5 45C34.9264 45 45 34.9264 45 22.5C45 10.0736 34.9264 0 22.5 0C10.0736 0 0 10.0736 0 22.5C0 34.9264 10.0736 45 22.5 45ZM9.71263 24.1166C10.3375 22.9715 11.5489 21.7848 12.9914 21.4657C14.193 21.1999 16.5523 21.3266 20.1139 24.8869C24.7053 29.4766 29.2944 31.0478 33.4664 30.125C37.3974 29.2555 39.9845 26.3673 41.2126 24.1166C42.1054 22.4804 41.5028 20.4302 39.8666 19.5374C38.2303 18.6445 36.1801 19.2472 35.2873 20.8834C34.6625 22.0285 33.4511 23.2152 32.0086 23.5343C30.807 23.8001 28.4476 23.6734 24.886 20.1131C20.2947 15.5234 15.7055 13.9522 11.5336 14.875C7.60257 15.7445 5.01544 18.6327 3.78733 20.8834C2.89452 22.5196 3.49717 24.5698 5.13339 25.4626C6.76962 26.3555 8.81981 25.7528 9.71263 24.1166Z"
        fill="#F4CE83"
      />
    </svg>
  );
}

export function ErrorIcon({ size }: StatusIconSizeProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 45 45"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M22.5 45C34.9264 45 45 34.9264 45 22.5C45 10.0736 34.9264 0 22.5 0C10.0736 0 0 10.0736 0 22.5C0 34.9264 10.0736 45 22.5 45ZM27 9.5C27 11.9853 26.25 31.25 22.5 31.25C18.75 31.25 18 11.9853 18 9.5C18 7.01472 20.0147 5 22.5 5C24.9853 5 27 7.01472 27 9.5ZM26.25 37.25C26.25 39.3211 24.5711 41 22.5 41C20.4289 41 18.75 39.3211 18.75 37.25C18.75 35.1789 20.4289 33.5 22.5 33.5C24.5711 33.5 26.25 35.1789 26.25 37.25Z"
        fill="#F48A83"
      />
    </svg>
  );
}

function getIconComponent(kind: StatusIconKind) {
  switch (kind) {
    case 'Loading':
      return LoadingIcon;
    case 'OK':
      return OKIcon;
    case 'Warning':
      return WarningIcon;
    case 'Error':
      return ErrorIcon;
  }
}

export type StatusIconKind = 'Loading' | 'OK' | 'Error' | 'Warning';

export type StatusIconProps = {
  kind: StatusIconKind;
} & Partial<StatusIconSizeProps> &
  ComponentPropsWithoutRef<'div'>;

export default function StatusIcon({ kind, size, ...other }: StatusIconProps) {
  size = size ?? defaultIconSize;
  const iconComponent = getIconComponent(kind);
  return <div {...other}>{iconComponent({ size })}</div>;
}
