/** 基于 [math-curve-loaders](https://github.com/Paidax01/math-curve-loaders) 的曲线配置 */
export type MathCurveVariant = 'thinking' | 'rose' | 'lissajous'

export interface MathCurvePoint {
  x: number
  y: number
}

export interface MathCurveConfig {
  rotate?: boolean
  particleCount: number
  trailSpan: number
  durationMs: number
  rotationDurationMs: number
  pulseDurationMs: number
  strokeWidth: number
  point: (progress: number, detailScale: number, config: MathCurveConfig) => MathCurvePoint
}

export interface MathCurveLoaderSize {
  xs: number
  sm: number
  md: number
  lg: number
  xl: number
}

export const MATH_CURVE_SIZES: MathCurveLoaderSize = {
  xs: 20,
  sm: 32,
  md: 48,
  lg: 72,
  xl: 96,
}
