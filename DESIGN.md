---
name: MUSEUM
colors:
  primary: '#191919'
  secondary: '#74716b'
  neutral: '#f4f1ea'
  border: '#d8d4cc'
  error: '#a63d32'
typography:
  heading:
    fontFamily: Georgia, 'Times New Roman', serif
    fontSize: 3rem
    fontWeight: 400
  body:
    fontFamily: Arial, Helvetica, sans-serif
    fontSize: 1rem
  caption:
    fontFamily: Arial, Helvetica, sans-serif
    fontSize: 0.9rem
    lineHeight: 1.8
rounded:
  none: 0px
spacing:
  xs: 0.5rem
  sm: 1rem
  md: 2rem
  lg: 4rem
  xl: 8rem
  2xl: 12rem
---

## Overview

MUSEUM is a quiet contemporary gallery: 自分という人間を編集して残す場所。
The hierarchy is Museum → Exhibition → Exhibit. Use 「展示する」 for the action.
Works, whitespace, typography and placement establish hierarchy.

## Colors

Use a warm neutral surface (#f4f1ea), dark ink (#191919), muted descriptions
(#74716b), fine borders (#d8d4cc), and restrained errors (#a63d32).
These values correspond to the existing --color-* variables in frontend/app/globals.css.

## Typography

Use the existing serif stack for room headings and work titles, and the sans-serif
stack for controls and captions. Headings have regular weight and catalogue-like
scale. Preserve the existing responsive room title sizing.

## Layout

Content width is 1200px; page padding is clamp(1.5rem, 5vw, 5rem).
Spacing tokens map to --space-xs through --space-2xl in globals.css.
Exhibits retain API display_order. Alternate left and right placement with 72% and
62% widths and 8rem separation; collapse to full width below 768px.
Display original image proportions, with titles and captions beneath the image.
The creation form sits after the works, at a maximum width of 640px.

## Elevation & Depth

Use whitespace and fine rules. Exhibit images have no card surface or shadow.

## Shapes

Use square corners. Do not add rounded containers around works.

## Components

Exhibition details retain their existing heading and description.
Exhibit images use browser-native rendering for arbitrary external hosts, with
lazy loading and descriptive alternative text. Forms use visible labels, fine
underlines, disabled submission while saving, and accessible status/error text.
An empty room and a failed request have distinct messages. A saved work remains
saved even when reloading fails; offer a list retry instead of another submission.

## Do's and Don'ts

- Use asymmetric catalogue placement and large whitespace.
- Preserve artwork proportions and user-authored captions.
- Keep 「新しい展示」 and 「展示する」 as the creation vocabulary.
- Avoid dashboards, feeds, uniform three-column cards, gradients, large CTAs,
  decorative shadows, and like/follower controls.
- Limit this change to Exhibit UI and documentation of existing design tokens.
