using System;

namespace Tachograph
{
    /// <summary>
    /// Represents the type of raw tachograph file.
    /// </summary>
    public enum RawFileType
    {
        Unknown = 0,
        Card = 1,
        VehicleUnit = 2
    }

    /// <summary>
    /// Represents a raw, unparsed tachograph file suitable for authentication.
    /// </summary>
    public class RawFile
    {
        /// <summary>
        /// Gets or sets the type of the raw file.
        /// </summary>
        public RawFileType Type { get; set; }

        /// <summary>
        /// Gets or sets the raw card file data (when Type is Card).
        /// </summary>
        public RawCardFile Card { get; set; }

        /// <summary>
        /// Gets or sets the raw vehicle unit file data (when Type is VehicleUnit).
        /// </summary>
        public RawVehicleUnitFile VehicleUnit { get; set; }

        /// <summary>
        /// Creates a deep clone of this RawFile.
        /// </summary>
        /// <returns>A cloned RawFile</returns>
        public RawFile Clone()
        {
            return new RawFile
            {
                Type = this.Type,
                Card = this.Card?.Clone(),
                VehicleUnit = this.VehicleUnit?.Clone()
            };
        }
    }
}
