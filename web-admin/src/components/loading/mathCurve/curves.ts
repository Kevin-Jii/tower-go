import type { MathCurveConfig, MathCurveVariant } from './types'

function rosePoint(
  progress: number,
  detailScale: number,
  cfg: MathCurveConfig & { baseRadius: number; detailAmplitude: number; petalCount: number; curveScale: number },
): { x: number; y: number } {
  const t = progress * Math.PI * 2
  const petals = Math.round(cfg.petalCount)
  const x = cfg.baseRadius * Math.cos(t) - cfg.detailAmplitude * detailScale * Math.cos(petals * t)
  const y = cfg.baseRadius * Math.sin(t) - cfg.detailAmplitude * detailScale * Math.sin(petals * t)
  return { x: 50 + x * cfg.curveScale, y: 50 + y * cfg.curveScale }
}

const thinking: MathCurveConfig & {
  baseRadius: number
  detailAmplitude: number
  petalCount: number
  curveScale: number
} = {
  rotate: true,
  particleCount: 64,
  trailSpan: 0.38,
  durationMs: 4600,
  rotationDurationMs: 28000,
  pulseDurationMs: 4200,
  strokeWidth: 5.5,
  baseRadius: 7,
  detailAmplitude: 3,
  petalCount: 7,
  curveScale: 3.9,
  point(progress, detailScale, config) {
    return rosePoint(progress, detailScale, config as typeof thinking)
  },
}

const rose: MathCurveConfig & { orbitRadius: number; petalCount: number; curveScale: number } = {
  rotate: true,
  particleCount: 58,
  trailSpan: 0.36,
  durationMs: 4800,
  rotationDurationMs: 32000,
  pulseDurationMs: 4400,
  strokeWidth: 5.2,
  orbitRadius: 7,
  petalCount: 7,
  curveScale: 3.6,
  point(progress, detailScale, config) {
    const c = config as typeof rose
    const t = progress * Math.PI * 2
    const radius = c.orbitRadius * (0.55 + detailScale * 0.45) * Math.cos(c.petalCount * t)
    return {
      x: 50 + radius * Math.cos(t) * c.curveScale,
      y: 50 + radius * Math.sin(t) * c.curveScale,
    }
  },
}

const lissajous: MathCurveConfig & { ax: number; ay: number; delta: number; curveScale: number } = {
  rotate: false,
  particleCount: 56,
  trailSpan: 0.34,
  durationMs: 5200,
  rotationDurationMs: 30000,
  pulseDurationMs: 4000,
  strokeWidth: 5,
  ax: 3,
  ay: 4,
  delta: Math.PI / 2,
  curveScale: 3.2,
  point(progress, detailScale, config) {
    const c = config as typeof lissajous
    const t = progress * Math.PI * 2
    const scale = c.curveScale * (0.75 + detailScale * 0.25)
    return {
      x: 50 + scale * 8 * Math.sin(c.ax * t + c.delta),
      y: 50 + scale * 8 * Math.sin(c.ay * t),
    }
  },
}

export const MATH_CURVE_PRESETS: Record<MathCurveVariant, MathCurveConfig> = {
  thinking,
  rose,
  lissajous,
}
