import type { Site } from "./site"
import type { ValueStream } from "./valueStream"

export interface Device {
  id: number
  name: string
  siteId?: number
  site?: Site
  valueStreamId?: number
  valueStream?: ValueStream
}