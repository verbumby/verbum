// Touch devices show an on-screen keyboard whenever an input gets focus, and it
// covers a good half of the screen. Auto-focusing is only helpful where a real
// keyboard is attached, which is what a fine pointer indicates. Narrow browser
// windows on computers must keep the desktop behaviour, so the viewport width
// deliberately plays no role here.
export const hasPhysicalKeyboard = (): boolean =>
    window.matchMedia('(any-pointer: fine)').matches
