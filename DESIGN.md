---
name: Job Market Analytics
description: Russian-language dashboard for exploring IT vacancies.
colors:
  bg: "#202221"
  surface: "#292b29"
  line: "#444640"
  ink: "#f2efe4"
  muted: "#bbbdb2"
  paper: "#e9e5d6"
  dark: "#292d26"
typography:
  body:
    fontFamily: '"Segoe UI", "Helvetica Neue", sans-serif'
    fontSize: "16px"
    lineHeight: 1.5
  display:
    fontSize: "clamp(36px, 4vw, 58px)"
    fontWeight: 500
    lineHeight: 1.14
    letterSpacing: "-0.04em"
  headline:
    fontSize: "clamp(19px, 1.7vw, 25px)"
    fontWeight: 550
    lineHeight: 1.25
    letterSpacing: "-0.03em"
rounded:
  control: "12px"
  panel: "22px"
  pill: "24px"
spacing:
  panel: "clamp(20px, 2.3vw, 36px)"
  grid: "20px"
components:
  button:
    backgroundColor: "{colors.surface}"
    rounded: "{rounded.control}"
  field:
    backgroundColor: "{colors.bg}"
    rounded: "{rounded.control}"
  current:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.dark}"
    rounded: "{rounded.pill}"
  skill:
    padding: "5px 12px"
    rounded: "16px"
  publications:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.dark}"
    rounded: "{rounded.panel}"
  period:
    backgroundColor: "{colors.dark}"
    textColor: "{colors.ink}"
    rounded: "{rounded.pill}"
---

## Overview

**Creative North Star: "The Market Overview"**

The Russian-language market overview leads with four counts and a wide ivory publication chart. Supporting panels use graphite so the reader can locate the time series immediately.

**Key Characteristics:**

- Tabular figures align counts across rows.
- Ivory marks the active navigation item and the publication chart.

## Colors

Graphite forms the canvas. A lighter graphite surface holds comparisons and lists. Warm ivory carries data marks on dark panels. On the publication chart, dark olive-gray bars sit on ivory and the selected bar deepens. Muted text carries labels and units.

## Typography

Segoe UI serves Russian labels and digits throughout. Page headings use close tracking and medium weight. Metric values grow with the viewport. The chart total and ring percentage share a 42px size. Keep the date or period alongside the chart value.

## Layout

Center content within 1320px. Below the metric strip, give publication dynamics 1.95 times the width of the adjacent work-format panel. At 880px and below, the ratio becomes 1.65 and the gap narrows. At 700px and below, stack panels in reading order, retain two metric columns and show four equal navigation items above the content.

## Elevation & Depth

Use surface contrast and thin dividers to separate sections. The publication panel enters over 450ms and bars grow from their baseline over 550ms. Reduced-motion preferences suppress both animations and control transitions.

## Shapes

Use rounded panels, smaller rounded fields, pill-shaped period switches and circular arrow links. Ranking bars run along a thin rounded track. Vertical publication bars have softer upper corners.

## Components

- **Publication dynamics.** Show the selected period total above the bars. Hovering or tapping a day pairs its count with its date. The daily disclosure exposes the underlying values.
- **Work format.** An ivory arc measures the remote share. Keep both category counts below the ring and use a dash when the denominator is zero.
- **Rankings.** Align names left and counts right. Scale bar lengths against the highest count. When available, percentages use the overall total.
- **Navigation and controls.** Pair navigation icons with labels. Use an ivory fill for the current route and dark fill for the selected period. Keyboard focus uses a two-pixel outline with a three-pixel offset.
- **Feedback.** Use muted text for loading and empty results. A failed request shows an alert with a retry button.

## Do's and Don'ts

- Do keep dates and counts together when a publication bar is selected.
- Do provide daily values in the chart’s disclosure.
- Don't rely on bar length alone to communicate a ranking.
