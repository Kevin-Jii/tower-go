import type { MathCurveConfig } from './types'

const SVG_NS = 'http://www.w3.org/2000/svg'

export interface MathCurveInstance {
  destroy: () => void
}

function normalizeProgress(progress: number): number {
  return ((progress % 1) + 1) % 1
}

function buildPath(config: MathCurveConfig, detailScale: number, steps = 360): string {
  return Array.from({ length: steps + 1 }, (_, index) => {
    const point = config.point(index / steps, detailScale, config)
    return `${index === 0 ? 'M' : 'L'} ${point.x.toFixed(2)} ${point.y.toFixed(2)}`
  }).join(' ')
}

function getDetailScale(time: number, config: MathCurveConfig, phaseOffset: number): number {
  const pulseProgress = (time % config.pulseDurationMs) / config.pulseDurationMs
  const pulseAngle = pulseProgress * Math.PI * 2
  return 0.52 + ((Math.sin(pulseAngle + phaseOffset * 0.55) + 1) / 2) * 0.48
}

function getRotation(time: number, config: MathCurveConfig, phaseOffset: number): number {
  if (!config.rotate) return 0
  return -(((time + phaseOffset * 1200) % config.rotationDurationMs) / config.rotationDurationMs) * 360
}

function getParticle(
  config: MathCurveConfig,
  index: number,
  progress: number,
  detailScale: number,
): { x: number; y: number; radius: number; opacity: number } {
  const tailOffset = index / (config.particleCount - 1)
  const point = config.point(
    normalizeProgress(progress - tailOffset * config.trailSpan),
    detailScale,
    config,
  )
  const fade = Math.pow(1 - tailOffset, 0.56)
  return {
    x: point.x,
    y: point.y,
    radius: 0.9 + fade * 2.7,
    opacity: 0.08 + fade * 0.92,
  }
}

/** 在容器内挂载 math-curve 粒子动画，返回销毁函数 */
export function mountMathCurveLoader(
  container: HTMLElement,
  config: MathCurveConfig,
): MathCurveInstance {
  const svg = document.createElementNS(SVG_NS, 'svg')
  svg.setAttribute('viewBox', '0 0 100 100')
  svg.setAttribute('fill', 'none')
  svg.setAttribute('aria-hidden', 'true')
  svg.setAttribute('class', 'math-curve-loader__svg')

  const group = document.createElementNS(SVG_NS, 'g')
  const path = document.createElementNS(SVG_NS, 'path')
  path.setAttribute('stroke', 'currentColor')
  path.setAttribute('stroke-width', String(config.strokeWidth))
  path.setAttribute('stroke-linecap', 'round')
  path.setAttribute('stroke-linejoin', 'round')
  path.setAttribute('opacity', '0.12')
  group.appendChild(path)
  svg.appendChild(group)

  const particles = Array.from({ length: config.particleCount }, () => {
    const circle = document.createElementNS(SVG_NS, 'circle')
    circle.setAttribute('fill', 'currentColor')
    group.appendChild(circle)
    return circle
  })

  container.replaceChildren(svg)

  const startTime = performance.now()
  const phaseOffset = Math.random()
  let frameId = 0
  let stopped = false

  function render(now: number): void {
    if (stopped) return
    const time = now - startTime
    const progress = ((time + phaseOffset * config.durationMs) % config.durationMs) / config.durationMs
    const detailScale = getDetailScale(time, config, phaseOffset)
    const rotation = getRotation(time, config, phaseOffset)

    group.setAttribute('transform', `rotate(${rotation} 50 50)`)
    path.setAttribute('d', buildPath(config, detailScale))

    particles.forEach((node, index) => {
      const particle = getParticle(config, index, progress, detailScale)
      node.setAttribute('cx', particle.x.toFixed(2))
      node.setAttribute('cy', particle.y.toFixed(2))
      node.setAttribute('r', particle.radius.toFixed(2))
      node.setAttribute('opacity', particle.opacity.toFixed(3))
    })

    frameId = requestAnimationFrame(render)
  }

  frameId = requestAnimationFrame(render)

  return {
    destroy() {
      stopped = true
      cancelAnimationFrame(frameId)
      container.replaceChildren()
    },
  }
}
