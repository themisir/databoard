import React, { useCallback, useState } from "react";

const App = () => {
  const [count, setCount] = useState(0);
  const increment = useCallback(() => {
    setCount((count) => count + 1);
  }, [count]);

  return (
    <>
      <h1>Hello!</h1>
      <h2>Count: {count}</h2>
      <button onClick={increment}>Increment</button>
    </>
  );
};

export default App;
