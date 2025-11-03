# app/schemas/__init__.py
import pkgutil
import importlib
import sys
from pathlib import Path

package = sys.modules[__name__]
package_path = Path(__file__).parent

for _, module_name, _ in pkgutil.iter_modules([str(package_path)]):
    module = importlib.import_module(f"{__name__}.{module_name}")
    for name in dir(module):
        if not name.startswith("_"):
            setattr(package, name, getattr(module, name))
