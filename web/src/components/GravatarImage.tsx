import { FC, RefAttributes, useCallback } from 'react'
import { Image } from 'react-bootstrap'
import { ImageProps } from 'react-bootstrap/Image'

export interface GravatarImageProps extends ImageProps, RefAttributes<HTMLImageElement> {
  profileUrl?: string
  baseSize: number
}

export const GravatarImage: FC<GravatarImageProps> = ({ profileUrl, baseSize, ...props }) => {
  const buildUrl = useCallback(
    (size: number) => {
      const url = new URL(profileUrl as string)
      url.searchParams.set('s', String(size))
      return url
    },
    [profileUrl],
  )

  return (
    profileUrl && (
      <Image
        src={buildUrl(baseSize).href}
        roundedCircle
        srcSet={[1, 1.5, 2, 2.5, 3, 3.5, 4].map((v) => buildUrl(v * baseSize).href + ' ' + v + 'x').join(',')}
        {...props}
      />
    )
  )
}

export default GravatarImage
