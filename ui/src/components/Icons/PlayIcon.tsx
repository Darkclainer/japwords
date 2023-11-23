import { forwardRef } from 'react';

import { IconProps } from './IconProps';

export default forwardRef<SVGSVGElement, IconProps>(function PlayIcon(
  { color = 'white', ...rest },
  forwardedRef,
) {
  return (
    <svg
      width="15"
      height="15"
      viewBox="10 10 41 41"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...rest}
      ref={forwardedRef}
    >
      <circle cx="31" cy="31" r="20" />
      <path
        d="M26.5904 21.7116C26.7411 21.6368 26.9118 21.6579 27.0393 21.7559L41.0688 30.7172C41.1679 30.7994 41.2311 30.9185 41.2311 31.0471C41.2311 31.1735 41.1669 31.2968 41.0688 31.378L27.0383 40.3382C26.9645 40.3983 26.8675 40.4299 26.7748 40.4299L26.5893 40.3878C26.446 40.3193 26.3553 40.1696 26.3553 40.0073V22.09C26.3553 21.9298 26.446 21.7791 26.5904 21.7116Z"
        fill={color}
      />
    </svg>
  );
});
