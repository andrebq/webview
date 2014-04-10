wrapGenerator.mark(gen);

function gen(x) {
  return wrapGenerator(function gen$($ctx0) {
    while (1) switch ($ctx0.prev = $ctx0.next) {
    case 0:
      $ctx0.next = 2;
      return x;
    case 2:
    case "end":
      return $ctx0.stop();
    }
  }, this);
}
