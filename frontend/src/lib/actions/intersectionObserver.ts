/**
 * A Svelte action that dispatches an 'intersect' event when the node enters the viewport.
 */
export function intersectionObserver(node: HTMLElement) {
  const observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting) {
      node.dispatchEvent(new CustomEvent('intersect'));
    }
  });

  observer.observe(node);

  return {
    destroy() {
      observer.unobserve(node);
    }
  };
}
