import importlib
import pkgutil

import app.models


def import_all_models():
    """Dynamically import all modules in app.models"""
    package = app.models
    for _, module_name, _ in pkgutil.iter_modules(package.__path__):
        importlib.import_module(f"{package.__name__}.{module_name}")
