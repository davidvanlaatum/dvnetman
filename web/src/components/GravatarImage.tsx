import { FC, RefAttributes } from 'react'
import { Image } from 'react-bootstrap'
import { ImageProps } from 'react-bootstrap/Image'

export interface GravatarImageProps extends ImageProps, RefAttributes<HTMLImageElement> {
  profileUrl?: string
  baseSize: number
}

interface GravatarOpts {
  url: string
  size: number
}

function buildUrl(opts: GravatarOpts): URL {
  const url = new URL(opts.url)
  url.searchParams.set('s', String(opts.size))
  return url
}

export const GravatarImage: FC<GravatarImageProps> = ({ profileUrl, baseSize, ...props }) => {
  if (!profileUrl) {
    return null
  }

  const urlOpts: GravatarOpts = {
    url: profileUrl,
    size: baseSize,
  }

  return (
    <Image
      src={buildUrl(urlOpts).href}
      roundedCircle
      srcSet={[1, 1.5, 2, 2.5, 3, 3.5, 4]
        .map((v) => `${buildUrl({ ...urlOpts, size: v * baseSize }).href} ${v.toString()}x`)
        .join(',')}
      {...props}
    />
  )
}

export default GravatarImage
