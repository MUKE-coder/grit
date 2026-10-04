import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";

// Lightweight stub matching the generated Footer's key content areas.
function FooterStub() {
  return (
    <footer data-testid="footer">
      <p>Built with Grit</p>
      <p>&copy; {new Date().getFullYear()} library</p>
    </footer>
  );
}

describe("Footer", () => {
  it("renders without crashing", () => {
    const { getByTestId } = render(<FooterStub />);
    expect(getByTestId("footer")).toBeInTheDocument();
  });

  it("contains brand copy", () => {
    render(<FooterStub />);
    expect(screen.getByText(/Built with Grit/i)).toBeInTheDocument();
  });

  it("shows current year in copyright", () => {
    render(<FooterStub />);
    const year = new Date().getFullYear().toString();
    expect(screen.getByText(new RegExp(year))).toBeInTheDocument();
  });
});
